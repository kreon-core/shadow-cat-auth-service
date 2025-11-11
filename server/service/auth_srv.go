package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"scs-auth-service/config"
	"scs-auth-service/helpers"
	"scs-auth-service/models/dto"
	"scs-auth-service/models/entity"
	"scs-auth-service/models/request"
)

type Auth struct {
	Config *config.Config
	AuthDB *gorm.DB
}

func NewAuthService(cfg *config.Config, authDB *gorm.DB) *Auth {
	return &Auth{
		Config: cfg,
		AuthDB: authDB,
	}
}

func (srv *Auth) Register(ctx context.Context, db *gorm.DB, req *request.RegisterReq) (*dto.AuthData, int, error) {
	if db == nil {
		db = srv.AuthDB
	}

	var authData dto.AuthData
	eCode := helpers.EDatabaseError
	err := db.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(ctx)

		user := entity.User{
			Username:     req.Username,
			Email:        req.Email,
			PasswordHash: fmt.Sprintf("%x", sha256.Sum256([]byte(req.Password))),
			DisplayName:  req.DisplayName,
			AvatarURL:    req.AvatarURL,
			Status:       "active",
		}
		txErr := tx.Create(&user).Error
		if txErr != nil {
			if errors.Is(txErr, gorm.ErrDuplicatedKey) {
				eCode = helpers.EUserAlreadyExists
				return fmt.Errorf("user already exists -> %w", txErr)
			}
			return fmt.Errorf("create user -> %w", txErr)
		}

		authData.UserID = user.ID.String()
		authData.Username = user.Username
		authData.Email = user.Email
		authData.Role = user.Role
		authData.DisplayName = user.DisplayName
		authData.AvatarURL = user.AvatarURL
		authData.Status = user.Status
		authData.PlayerID = user.PlayerID.String()

		jwtIssuer := srv.Config.Secrets.JWTIssuer
		jwtSecretKey := []byte(srv.Config.Secrets.JWTSecretKey)

		jwtAccessToken, txErr := helpers.GenerateJWTAccessToken(
			user.ID.String(), user.PlayerID.String(), user.Role,
			jwtIssuer, jwtSecretKey,
			srv.Config.Secrets.AccessTokenExpiry,
		)
		if txErr != nil {
			eCode = helpers.EJWTGenerationFailed
			return fmt.Errorf("generate jwt access token -> %w", txErr)
		}

		jwtRefreshToken, token, txErr := helpers.GenerateJWTRefreshToken(
			user.ID.String(), user.PlayerID.String(), user.Role,
			jwtIssuer, jwtSecretKey,
			srv.Config.Secrets.RefreshTokenExpiry,
		)
		if txErr != nil {
			eCode = helpers.EJWTGenerationFailed
			return fmt.Errorf("generate jwt refresh token -> %w", txErr)
		}

		authData.TokenType = "Bearer"
		authData.AccessToken = jwtAccessToken
		authData.ExpiresIn = int64(srv.Config.Secrets.AccessTokenExpiry.Seconds())
		authData.RefreshToken = jwtRefreshToken

		userSession := entity.UserSession{
			UserID:    user.ID,
			Token:     token,
			ExpiresAt: time.Now().Add(srv.Config.Secrets.RefreshTokenExpiry),
			// DeviceInfo & IP Address
		}
		txErr = tx.Create(&userSession).Error
		if txErr != nil {
			if errors.Is(txErr, gorm.ErrDuplicatedKey) {
				eCode = helpers.EResourceAlreadyExists
				return fmt.Errorf("user session already exists -> %w", txErr)
			}
			return fmt.Errorf("create user session -> %w", txErr)
		}

		return nil
	})
	if err != nil {
		return nil, eCode, fmt.Errorf("register_transaction -> %w", err)
	}
	return &authData, helpers.Success, nil
}

func (srv *Auth) Login(ctx context.Context, db *gorm.DB, req *request.LoginReq) (*dto.AuthData, int, error) {
	if db == nil {
		db = srv.AuthDB
	}

	var authData dto.AuthData
	eCode := helpers.EDatabaseError
	err := db.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(ctx)

		var user entity.User
		txErr := tx.Where("username = ?", req.Username).First(&user).Error
		if txErr != nil {
			if errors.Is(txErr, gorm.ErrRecordNotFound) {
				eCode = helpers.EUserNotFound
				return errors.New("invalid credentials")
			}
			return fmt.Errorf("fetch user -> %w", txErr)
		}

		passwordHash := fmt.Sprintf("%x", sha256.Sum256([]byte(req.Password)))
		if user.PasswordHash != passwordHash {
			eCode = helpers.EAccessDenied
			return errors.New("invalid credentials")
		}

		authData.UserID = user.ID.String()
		authData.Username = user.Username
		authData.Email = user.Email
		authData.Role = user.Role
		authData.DisplayName = user.DisplayName
		authData.AvatarURL = user.AvatarURL
		authData.Status = user.Status
		authData.PlayerID = user.PlayerID.String()

		jwtIssuer := srv.Config.Secrets.JWTIssuer
		jwtSecretKey := []byte(srv.Config.Secrets.JWTSecretKey)

		jwtAccessToken, txErr := helpers.GenerateJWTAccessToken(
			user.ID.String(), user.PlayerID.String(), user.Role,
			jwtIssuer, jwtSecretKey,
			srv.Config.Secrets.AccessTokenExpiry,
		)
		if txErr != nil {
			eCode = helpers.EJWTGenerationFailed
			return fmt.Errorf("generate jwt access token -> %w", txErr)
		}

		jwtRefreshToken, token, txErr := helpers.GenerateJWTRefreshToken(
			user.ID.String(), user.PlayerID.String(), user.Role,
			jwtIssuer, jwtSecretKey,
			srv.Config.Secrets.RefreshTokenExpiry,
		)
		if txErr != nil {
			eCode = helpers.EJWTGenerationFailed
			return fmt.Errorf("generate jwt refresh token -> %w", txErr)
		}

		authData.TokenType = "Bearer"
		authData.AccessToken = jwtAccessToken
		authData.ExpiresIn = int64(srv.Config.Secrets.AccessTokenExpiry.Seconds())
		authData.RefreshToken = jwtRefreshToken

		userSession := entity.UserSession{
			UserID:    user.ID,
			Token:     token,
			ExpiresAt: time.Now().Add(srv.Config.Secrets.RefreshTokenExpiry),
			// DeviceInfo & IP Address
		}
		txErr = tx.Create(&userSession).Error
		if txErr != nil {
			if errors.Is(txErr, gorm.ErrDuplicatedKey) {
				eCode = helpers.EResourceAlreadyExists
				return fmt.Errorf("user session already exists -> %w", txErr)
			}
			return fmt.Errorf("create user session -> %w", txErr)
		}

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

func (srv *Auth) Logout(ctx context.Context, db *gorm.DB, req *request.LogoutReq) (int, error) {
	if db == nil {
		db = srv.AuthDB
	}

	eCode := helpers.EDatabaseError
	err := db.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(ctx)

		var user entity.User
		txErr := tx.Where("id = ?", req.UserID).First(&user).Error
		if txErr != nil {
			eCode = helpers.EUserNotFound
			return fmt.Errorf("fetch user -> %w", txErr)
		}

		txErr = tx.Model(&entity.UserSession{}).
			Where("user_id = ? AND token = ? AND revoked = false",
				user.ID, req.RefreshToken).
			Update("revoked", true).Error
		if txErr != nil {
			eCode = helpers.EResourceNotFound
			return fmt.Errorf("revoke user session -> %w", txErr)
		}

		return nil
	})
	if err != nil {
		return eCode, fmt.Errorf("logout_transaction -> %w", err)
	}
	return helpers.Success, nil
}
