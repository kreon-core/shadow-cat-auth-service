package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"scs-auth-service/config"
	"scs-auth-service/helpers"
	"scs-auth-service/models/dto"
	"scs-auth-service/models/entity"
	"scs-auth-service/models/enum"
	"scs-auth-service/models/request"
	"scs-auth-service/server/repository"
)

type Auth struct {
	Config   *config.Config
	AuthRepo *repository.Auth
}

func NewAuthService(cfg *config.Config, authRepo *repository.Auth) *Auth {
	return &Auth{
		Config:   cfg,
		AuthRepo: authRepo,
	}
}

func (srv *Auth) Register(ctx context.Context, req *request.RegisterReq) (*dto.AuthData, int, error) {
	authData := &dto.AuthData{}
	eCode := helpers.EDatabaseError

	err := srv.AuthRepo.WithTransaction(ctx, nil, func(tx *gorm.DB) error {
		passwordHash, err := helpers.HashUsingBcrypt(req.Password)
		if err != nil {
			eCode = helpers.EInvalidRequest
			return fmt.Errorf("hash_password -> %w", err)
		}

		emailLocal := strings.SplitN(req.Email, "@", 2)[0]
		username := helpers.OrDefault(req.Username, emailLocal)
		user := &entity.User{
			Username:     username,
			Email:        req.Email,
			PasswordHash: passwordHash,
			Role:         "player",
			DisplayName:  helpers.OrDefault(req.DisplayName, req.Username),
			AvatarURL: helpers.OrDefault(
				req.AvatarURL,
				fmt.Sprintf("https://www.gravatar.com/avatar/%s?d=identicon&s=128", username),
			),
			Status: "active",
		}
		txErr := srv.AuthRepo.CUser(ctx, tx, user)
		if txErr != nil {
			eCode = helpers.ConvertPgErrToAppCode(txErr)
			return fmt.Errorf("create_user -> %w", txErr)
		}

		jwtIssuer := srv.Config.Secrets.JWTIssuer
		jwtSecretKey := []byte(srv.Config.Secrets.JWTSecretKey)
		jwtAccessToken, txErr := helpers.GenerateJWTAccessToken(
			user.ID.String(), user.PlayerID.String(), user.Role,
			jwtIssuer, jwtSecretKey,
			srv.Config.Secrets.AccessTokenExpiry,
		)
		if txErr != nil {
			eCode = helpers.EJWTGenerationFailed
			return fmt.Errorf("generate_jwt_access_token -> %w", txErr)
		}
		jwtRefreshToken, token, txErr := helpers.GenerateJWTRefreshToken(
			user.ID.String(), user.PlayerID.String(), user.Role,
			jwtIssuer, jwtSecretKey,
			srv.Config.Secrets.RefreshTokenExpiry,
		)
		if txErr != nil {
			eCode = helpers.EJWTGenerationFailed
			return fmt.Errorf("generate_jwt_refresh_token -> %w", txErr)
		}
		userSession := &entity.UserSession{
			UserID:    user.ID,
			Token:     token,
			ExpiresAt: time.Now().Add(srv.Config.Secrets.RefreshTokenExpiry),
			// DeviceInfo & IP Address
		}
		txErr = srv.AuthRepo.SUserSession(ctx, tx, userSession)
		if txErr != nil {
			eCode = helpers.ConvertPgErrToAppCode(txErr)
			return fmt.Errorf("create_user_session -> %w", txErr)
		}

		authData.UserID = user.ID.String()
		authData.Username = user.Username
		authData.Email = user.Email
		authData.Role = user.Role
		authData.DisplayName = user.DisplayName
		authData.AvatarURL = user.AvatarURL
		authData.Status = user.Status
		authData.PlayerID = user.PlayerID.String()
		authData.TokenType = enum.TokenTypeBearer
		authData.AccessToken = jwtAccessToken
		authData.ExpiresIn = int64(srv.Config.Secrets.AccessTokenExpiry.Seconds())
		authData.RefreshToken = jwtRefreshToken

		return nil
	})
	if err != nil {
		return nil, eCode, fmt.Errorf("register_transaction -> %w", err)
	}
	return authData, helpers.Success, nil
}

func (srv *Auth) Login(ctx context.Context, req *request.LoginReq) (*dto.AuthData, int, error) {
	authData := &dto.AuthData{}
	eCode := helpers.EDatabaseError

	err := srv.AuthRepo.WithTransaction(ctx, nil, func(tx *gorm.DB) error {
		var user *entity.User
		var txErr error

		if !helpers.IsBlankString(&req.Username) {
			user, txErr = srv.AuthRepo.RUserWUsername(ctx, tx, req.Username)
			if txErr != nil {
				eCode = helpers.ConvertPgErrToAppCode(txErr)
				return fmt.Errorf("fetch_user -> %w", txErr)
			}
		} else if !helpers.IsBlankString(&req.Email) {
			user, txErr = srv.AuthRepo.RUserWEmail(ctx, tx, req.Email)
			if txErr != nil {
				eCode = helpers.ConvertPgErrToAppCode(txErr)
				return fmt.Errorf("fetch_user -> %w", txErr)
			}
		} else {
			eCode = helpers.EInvalidRequest
			return errors.New("missing_username_or_email")
		}

		if user.Status != "active" {
			eCode = helpers.EAccountSuspended
			return errors.New("user_inactive")
		}

		if !helpers.VerifyBcryptHash(user.PasswordHash, req.Password) {
			eCode = helpers.EAccessDenied
			return errors.New("invalid_credentials")
		}

		jwtIssuer := srv.Config.Secrets.JWTIssuer
		jwtSecretKey := []byte(srv.Config.Secrets.JWTSecretKey)
		jwtAccessToken, txErr := helpers.GenerateJWTAccessToken(
			user.ID.String(), user.PlayerID.String(), user.Role,
			jwtIssuer, jwtSecretKey,
			srv.Config.Secrets.AccessTokenExpiry,
		)
		if txErr != nil {
			eCode = helpers.EJWTGenerationFailed
			return fmt.Errorf("generate_jwt_access_token -> %w", txErr)
		}
		jwtRefreshToken, token, txErr := helpers.GenerateJWTRefreshToken(
			user.ID.String(), user.PlayerID.String(), user.Role,
			jwtIssuer, jwtSecretKey,
			srv.Config.Secrets.RefreshTokenExpiry,
		)
		if txErr != nil {
			eCode = helpers.EJWTGenerationFailed
			return fmt.Errorf("generate_jwt_refresh_token -> %w", txErr)
		}
		userSession := &entity.UserSession{
			UserID:    user.ID,
			Token:     token,
			ExpiresAt: time.Now().Add(srv.Config.Secrets.RefreshTokenExpiry),
			// DeviceInfo & IP Address
		}
		txErr = srv.AuthRepo.SUserSession(ctx, tx, userSession)
		if txErr != nil {
			eCode = helpers.ConvertPgErrToAppCode(txErr)
			return fmt.Errorf("create_user_session -> %w", txErr)
		}

		authData.UserID = user.ID.String()
		authData.Username = user.Username
		authData.Email = user.Email
		authData.Role = user.Role
		authData.DisplayName = user.DisplayName
		authData.AvatarURL = user.AvatarURL
		authData.Status = user.Status
		authData.PlayerID = user.PlayerID.String()
		authData.TokenType = enum.TokenTypeBearer
		authData.AccessToken = jwtAccessToken
		authData.ExpiresIn = int64(srv.Config.Secrets.AccessTokenExpiry.Seconds())
		authData.RefreshToken = jwtRefreshToken

		return nil
	})
	if err != nil {
		return nil, eCode, fmt.Errorf("login_transaction -> %w", err)
	}
	return authData, helpers.Success, nil
}

func (srv *Auth) AuthZalo(ctx context.Context, req *request.AuthZaloReq) (*dto.AuthData, int, error) {
	authData := &dto.AuthData{}
	eCode := helpers.EDatabaseError

	err := srv.AuthRepo.WithTransaction(ctx, nil, func(tx *gorm.DB) error {
		var user *entity.User
		var authProvider *entity.AuthProvider
		isNewUser := false

		authProvider, txErr := srv.AuthRepo.RAuthProviderWProviderUID(ctx, tx, "zalo", req.UID)
		if txErr != nil {
			if errors.Is(txErr, gorm.ErrRecordNotFound) {
				isNewUser = true
			} else {
				eCode = helpers.ConvertPgErrToAppCode(txErr)
				return fmt.Errorf("fetch_user_by_auth_provider -> %w", txErr)
			}
		}

		if !isNewUser && authProvider == nil {
			isNewUser = true
		}
		if isNewUser {
			user = &entity.User{
				Username:    fmt.Sprintf("zalo_%s", req.UID),
				Role:        "player",
				DisplayName: helpers.OrDefault(req.DisplayName, fmt.Sprintf("ZaloUser_%s", req.UID)),
				AvatarURL: helpers.OrDefault(
					req.AvatarURL,
					fmt.Sprintf("https://www.gravatar.com/avatar/%s?d=identicon&s=128", req.UID),
				),
				Status: "active",
			}
			txErr = srv.AuthRepo.CUser(ctx, tx, user)
			if txErr != nil {
				eCode = helpers.ConvertPgErrToAppCode(txErr)
				return fmt.Errorf("create_user -> %w", txErr)
			}

			var metadata datatypes.JSON
			metadataMap := map[string]any{
				"utm_source":   req.UtmSource,
				"utm_medium":   req.UtmMedium,
				"utm_campaign": req.UtmCampaign,
			}
			metadataJSON, jsonErr := json.Marshal(metadataMap)
			if jsonErr == nil {
				metadata = datatypes.JSON(metadataJSON)
			}
			authProvider = &entity.AuthProvider{
				UserID:      user.ID,
				Provider:    "zalo",
				ProviderUID: req.UID,
				LinkedAt:    time.Now(),
				Metadata:    metadata,
			}
			txErr = srv.AuthRepo.CAuthProvider(ctx, tx, authProvider)
			if txErr != nil {
				eCode = helpers.ConvertPgErrToAppCode(txErr)
				return fmt.Errorf("create_auth_provider -> %w", txErr)
			}
		} else {
			user, txErr = srv.AuthRepo.RUserWID(ctx, tx, authProvider.UserID)
			if txErr != nil {
				eCode = helpers.ConvertPgErrToAppCode(txErr)
				return fmt.Errorf("fetch_user -> %w", txErr)
			}
			if user.Status != "active" {
				eCode = helpers.EAccountSuspended
				return errors.New("user_inactive")
			}
		}

		jwtIssuer := srv.Config.Secrets.JWTIssuer
		jwtSecretKey := []byte(srv.Config.Secrets.JWTSecretKey)
		jwtAccessToken, txErr := helpers.GenerateJWTAccessToken(
			user.ID.String(), user.PlayerID.String(), user.Role,
			jwtIssuer, jwtSecretKey,
			srv.Config.Secrets.AccessTokenExpiry,
		)
		if txErr != nil {
			eCode = helpers.EJWTGenerationFailed
			return fmt.Errorf("generate_jwt_access_token -> %w", txErr)
		}
		jwtRefreshToken, token, txErr := helpers.GenerateJWTRefreshToken(
			user.ID.String(), user.PlayerID.String(), user.Role,
			jwtIssuer, jwtSecretKey,
			srv.Config.Secrets.RefreshTokenExpiry,
		)
		if txErr != nil {
			eCode = helpers.EJWTGenerationFailed
			return fmt.Errorf("generate_jwt_refresh_token -> %w", txErr)
		}
		userSession := &entity.UserSession{
			UserID:    user.ID,
			Token:     token,
			ExpiresAt: time.Now().Add(srv.Config.Secrets.RefreshTokenExpiry),
			// DeviceInfo & IP Address
		}
		txErr = srv.AuthRepo.SUserSession(ctx, tx, userSession)
		if txErr != nil {
			eCode = helpers.ConvertPgErrToAppCode(txErr)
			return fmt.Errorf("create_user_session -> %w", txErr)
		}

		authData.UserID = user.ID.String()
		authData.Username = user.Username
		authData.Email = user.Email
		authData.Role = user.Role
		authData.DisplayName = user.DisplayName
		authData.AvatarURL = user.AvatarURL
		authData.Status = user.Status
		authData.Provider = &dto.AuthProviderData{
			Name:     authProvider.Provider,
			UID:      authProvider.ProviderUID,
			LinkedAt: authProvider.LinkedAt,
			Metadata: authProvider.Metadata,
		}
		authData.PlayerID = user.PlayerID.String()
		authData.TokenType = enum.TokenTypeBearer
		authData.AccessToken = jwtAccessToken
		authData.ExpiresIn = int64(srv.Config.Secrets.AccessTokenExpiry.Seconds())
		authData.RefreshToken = jwtRefreshToken

		return nil
	})
	if err != nil {
		return nil, eCode, fmt.Errorf("auth_zalo_transaction -> %w", err)
	}
	return authData, helpers.Success, nil
}

func (srv *Auth) AuthFirebase() {}

func (srv *Auth) AuthRefresh(ctx context.Context, req *request.AuthRefreshReq) (*dto.AuthRefreshData, int, error) {
	authRefreshData := &dto.AuthRefreshData{}
	eCode := helpers.EDatabaseError

	err := srv.AuthRepo.WithTransaction(ctx, nil, func(tx *gorm.DB) error {
		userID, parseErr := uuid.Parse(req.UserID)
		if parseErr != nil {
			eCode = helpers.EInvalidUUIDFormat
			return fmt.Errorf("parse_user_id -> %w", parseErr)
		}

		user, txErr := srv.AuthRepo.RUserWID(ctx, tx, userID)
		if txErr != nil {
			eCode = helpers.ConvertPgErrToAppCode(txErr)
			return fmt.Errorf("fetch_user -> %w", txErr)
		}

		token, claims, txErr := helpers.ParseJWTRefreshToken(req.RefreshToken, []byte(srv.Config.Secrets.JWTSecretKey))
		if txErr != nil {
			if errors.Is(txErr, jwt.ErrTokenExpired) {
				eCode = helpers.EExpiredRefreshToken
			} else {
				eCode = helpers.EInvalidRefreshToken
			}
			return fmt.Errorf("parse_jwt_refresh_token -> %w", txErr)
		}
		if !token.Valid {
			eCode = helpers.EInvalidRefreshToken
			return errors.New("invalid_jwt_refresh_token")
		}

		if claims.UserID != user.ID.String() {
			eCode = helpers.EInvalidRefreshToken
			return errors.New("token_user_mismatch")
		}

		userSession, txErr := srv.AuthRepo.RUserSessionWUserIDAndToken(ctx, tx, user.ID, claims.Token)
		if txErr != nil {
			eCode = helpers.ConvertPgErrToAppCode(txErr)
			return fmt.Errorf("fetch_user_session -> %w", txErr)
		}

		if userSession.Revoked {
			eCode = helpers.EAccessDenied
			return errors.New("refresh_token_revoked")
		}

		if userSession.ExpiresAt.Before(time.Now()) {
			eCode = helpers.EExpiredRefreshToken
			return errors.New("refresh_token_expired")
		}

		jwtIssuer := srv.Config.Secrets.JWTIssuer
		jwtSecretKey := []byte(srv.Config.Secrets.JWTSecretKey)
		jwtAccessToken, txErr := helpers.GenerateJWTAccessToken(
			user.ID.String(), user.PlayerID.String(), user.Role,
			jwtIssuer, jwtSecretKey,
			srv.Config.Secrets.AccessTokenExpiry,
		)
		if txErr != nil {
			eCode = helpers.EJWTGenerationFailed
			return fmt.Errorf("generate_jwt_access_token -> %w", txErr)
		}
		jwtRefreshToken, newToken, txErr := helpers.GenerateJWTRefreshToken(
			user.ID.String(), user.PlayerID.String(), user.Role,
			jwtIssuer, jwtSecretKey,
			srv.Config.Secrets.RefreshTokenExpiry,
		)
		if txErr != nil {
			eCode = helpers.EJWTGenerationFailed
			return fmt.Errorf("generate_jwt_refresh_token -> %w", txErr)
		}

		userSession.Token = newToken
		userSession.ExpiresAt = time.Now().Add(srv.Config.Secrets.RefreshTokenExpiry)
		txErr = srv.AuthRepo.UUserSession(ctx, tx, userSession)
		if txErr != nil {
			eCode = helpers.ConvertPgErrToAppCode(txErr)
			return fmt.Errorf("update_user_session -> %w", txErr)
		}

		authRefreshData.TokenType = enum.TokenTypeBearer
		authRefreshData.AccessToken = jwtAccessToken
		authRefreshData.ExpiresIn = int64(srv.Config.Secrets.AccessTokenExpiry.Seconds())
		authRefreshData.RefreshToken = jwtRefreshToken

		return nil
	})
	if err != nil {
		return nil, eCode, fmt.Errorf("auth_refresh_transaction -> %w", err)
	}
	return authRefreshData, helpers.Success, nil
}

func (srv *Auth) Logout(ctx context.Context, userID string) (int, error) {
	eCode := helpers.EDatabaseError

	err := srv.AuthRepo.WithTransaction(ctx, nil, func(tx *gorm.DB) error {
		userID, parseErr := uuid.Parse(userID)
		if parseErr != nil {
			eCode = helpers.EInvalidUUIDFormat
			return fmt.Errorf("parse_user_id -> %w", parseErr)
		}

		user, txErr := srv.AuthRepo.RUserWID(ctx, tx, userID)
		if txErr != nil {
			eCode = helpers.ConvertPgErrToAppCode(txErr)
			return fmt.Errorf("fetch_user -> %w", txErr)
		}

		txErr = srv.AuthRepo.UUserSessionFuncRevokeSession(ctx, tx, user.ID)
		if txErr != nil {
			eCode = helpers.ConvertPgErrToAppCode(txErr)
			return fmt.Errorf("revoke_user_session -> %w", txErr)
		}

		return nil
	})
	if err != nil {
		return eCode, fmt.Errorf("logout_transaction -> %w", err)
	}
	return helpers.Success, nil
}
