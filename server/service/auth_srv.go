package service

import (
	"gorm.io/gorm"

	"scs-auth-service/config"
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
