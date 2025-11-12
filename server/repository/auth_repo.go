package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"scs-auth-service/config"
	"scs-auth-service/models/entity"
)

type Auth struct {
	Config *config.Config
	AuthDB *gorm.DB
}

func NewAuthRepository(config *config.Config, authDB *gorm.DB) *Auth {
	return &Auth{
		Config: config,
		AuthDB: authDB,
	}
}

func (repo *Auth) GetDB() *gorm.DB {
	return repo.AuthDB
}

func (repo *Auth) WithTransaction(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error) error {
	if db == nil {
		db = repo.AuthDB
	}
	return db.WithContext(ctx).Transaction(fn)
}

func (repo *Auth) CUser(ctx context.Context, db *gorm.DB, user *entity.User) error {
	if db == nil {
		db = repo.AuthDB
	}
	return db.WithContext(ctx).Clauses(clause.Returning{}).Create(user).Error
}

func (repo *Auth) RUserWID(ctx context.Context, db *gorm.DB, userID string) (*entity.User, error) {
	if db == nil {
		db = repo.AuthDB
	}
	var user entity.User
	err := db.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *Auth) RUserWUsername(ctx context.Context, db *gorm.DB, username string) (*entity.User, error) {
	if db == nil {
		db = repo.AuthDB
	}
	var user entity.User
	err := db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *Auth) CUserSession(ctx context.Context, db *gorm.DB, session *entity.UserSession) error {
	if db == nil {
		db = repo.AuthDB
	}
	return db.WithContext(ctx).Clauses(clause.Returning{}).Create(session).Error
}

func (repo *Auth) RUserSessionWUserIDAndToken(
	ctx context.Context,
	db *gorm.DB,
	userID uuid.UUID,
	token string,
) (*entity.UserSession, error) {
	if db == nil {
		db = repo.AuthDB
	}
	var session entity.UserSession
	err := db.WithContext(ctx).Where("user_id = ? AND token = ?", userID, token).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (repo *Auth) UUserSession(ctx context.Context, db *gorm.DB, session *entity.UserSession) error {
	if db == nil {
		db = repo.AuthDB
	}
	return db.WithContext(ctx).Save(session).Error
}

func (repo *Auth) UUserSessionFuncRevokeSession(
	ctx context.Context,
	db *gorm.DB,
	userID uuid.UUID,
	token string,
) error {
	if db == nil {
		db = repo.AuthDB
	}
	err := db.WithContext(ctx).Model(&entity.UserSession{}).
		Where("user_id = ? AND token = ? AND revoked = false", userID, token).
		Update("revoked", true).Error
	return err
}

func (repo *Auth) CAuthProvider(ctx context.Context, db *gorm.DB, provider *entity.AuthProvider) error {
	if db == nil {
		db = repo.AuthDB
	}
	return db.WithContext(ctx).Clauses(clause.Returning{}).Create(provider).Error
}

func (repo *Auth) RAuthProviderWProviderUID(
	ctx context.Context,
	db *gorm.DB,
	provider, uid string,
) (*entity.AuthProvider, error) {
	if db == nil {
		db = repo.AuthDB
	}
	var authProvider entity.AuthProvider
	err := db.WithContext(ctx).Where("provider = ? AND uid = ?", provider, uid).First(&authProvider).Error
	if err != nil {
		return nil, err
	}
	return &authProvider, nil
}
