package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"scs-auth-service/config"
	"scs-auth-service/helpers"
	"scs-auth-service/models/dto"
	"scs-auth-service/models/entity"
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
	var authData dto.AuthData
	eCode := helpers.EDatabaseError

	err := srv.AuthRepo.WithTransaction(ctx, nil, func(tx *gorm.DB) error {
		passwordHash, err := helpers.HashUsingBcrypt(req.Password)
		if err != nil {
			eCode = helpers.EInvalidRequest
			return fmt.Errorf("hash_password -> %w", err)
		}
		user := entity.User{
			Username:     req.Username,
			Email:        req.Email,
			PasswordHash: passwordHash,
			DisplayName:  helpers.OrDefault(req.DisplayName, req.Username),
			AvatarURL: helpers.OrDefault(
				req.AvatarURL,
				fmt.Sprintf("https://www.gravatar.com/avatar/%s?d=identicon&s=128", req.Username),
			),
			Status: "active",
		}
		txErr := srv.AuthRepo.CUser(ctx, tx, &user)
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
		userSession := entity.UserSession{
			UserID:    user.ID,
			Token:     token,
			ExpiresAt: time.Now().Add(srv.Config.Secrets.RefreshTokenExpiry),
			// DeviceInfo & IP Address
		}
		txErr = srv.AuthRepo.CUserSession(ctx, tx, &userSession)
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
		authData.TokenType = "Bearer"
		authData.AccessToken = jwtAccessToken
		authData.ExpiresIn = int64(srv.Config.Secrets.AccessTokenExpiry.Seconds())
		authData.RefreshToken = jwtRefreshToken

		return nil
	})
	if err != nil {
		return nil, eCode, fmt.Errorf("register_transaction -> %w", err)
	}
	return &authData, helpers.Success, nil
}

func (srv *Auth) Login(ctx context.Context, req *request.LoginReq) (*dto.AuthData, int, error) {
	var authData dto.AuthData
	eCode := helpers.EDatabaseError

	err := srv.AuthRepo.WithTransaction(ctx, nil, func(tx *gorm.DB) error {
		user, txErr := srv.AuthRepo.RUserWUsername(ctx, tx, req.Username)
		if txErr != nil {
			eCode = helpers.ConvertPgErrToAppCode(txErr)
			return fmt.Errorf("fetch_user -> %w", txErr)
		}

		if user.Status != "active" {
			eCode = helpers.EAccountSuspended
			return errors.New("user_inactive")
		}

		if helpers.VerifyBcryptHash(user.PasswordHash, req.Password) {
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
		userSession := entity.UserSession{
			UserID:    user.ID,
			Token:     token,
			ExpiresAt: time.Now().Add(srv.Config.Secrets.RefreshTokenExpiry),
			// DeviceInfo & IP Address
		}
		txErr = srv.AuthRepo.CUserSession(ctx, tx, &userSession)
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
		authData.TokenType = "Bearer"
		authData.AccessToken = jwtAccessToken
		authData.ExpiresIn = int64(srv.Config.Secrets.AccessTokenExpiry.Seconds())
		authData.RefreshToken = jwtRefreshToken

		return nil
	})
	if err != nil {
		return nil, eCode, fmt.Errorf("login_transaction -> %w", err)
	}
	return &authData, helpers.Success, nil
}

func (srv *Auth) AuthZalo()     {}
func (srv *Auth) AuthFirebase() {}
func (srv *Auth) RefreshToken() {}

func (srv *Auth) Logout(ctx context.Context, req *request.LogoutReq) (int, error) {
	eCode := helpers.EDatabaseError

	err := srv.AuthRepo.WithTransaction(ctx, nil, func(tx *gorm.DB) error {
		user, txErr := srv.AuthRepo.RUserWID(ctx, tx, req.UserID)
		if txErr != nil {
			eCode = helpers.ConvertPgErrToAppCode(txErr)
			return fmt.Errorf("fetch_user -> %w", txErr)
		}

		token, claims, txErr := helpers.ParseJWTRefreshToken(req.RefreshToken, []byte(srv.Config.Secrets.JWTSecretKey))
		if txErr != nil {
			if !errors.Is(txErr, jwt.ErrTokenExpired) {
				eCode = helpers.EInvalidRefreshToken
				return fmt.Errorf("parse_jwt_refresh_token -> %w", txErr)
			}
		}
		if !token.Valid {
			eCode = helpers.EInvalidRefreshToken
			return errors.New("invalid_jwt_refresh_token")
		}

		if claims.UserID != user.ID.String() {
			eCode = helpers.EInvalidRefreshToken
			return errors.New("token_user_mismatch")
		}

		txErr = srv.AuthRepo.UUserSessionFuncRevokeSession(ctx, tx, user.ID, claims.Token)
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
