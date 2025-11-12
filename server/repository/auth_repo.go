package repository

import (
	"context"

	"gorm.io/gorm"

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

func (repo *Auth) CreateUser(ctx context.Context, db *gorm.DB, user *entity.User) error {
	if db == nil {
		db = repo.AuthDB
	}
	return db.WithContext(ctx).Create(user).Error
}

func (repo *Auth) CreateUserSession(ctx context.Context, db *gorm.DB, session *entity.UserSession) error {
	if db == nil {
		db = repo.AuthDB
	}
	return db.WithContext(ctx).Create(session).Error
}

func (repo *Auth) CreateAuthProvider(ctx context.Context, db *gorm.DB, provider *entity.AuthProvider) error {
	if db == nil {
		db = repo.AuthDB
	}
	return db.WithContext(ctx).Create(provider).Error
}
