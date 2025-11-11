package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"scs-auth-service/config"
	"scs-auth-service/server/service"
)

type Auth struct {
	Config  *config.Config
	AuthSrv *service.Auth
}

func NewAuthController(cfg *config.Config, authDB *gorm.DB) *Auth {
	return &Auth{
		Config:  cfg,
		AuthSrv: service.NewAuthService(cfg, authDB),
	}
}

func (ctrl *Auth) Register(c *gin.Context)     {}
func (ctrl *Auth) Login(c *gin.Context)        {}
func (ctrl *Auth) AuthZalo(c *gin.Context)     {}
func (ctrl *Auth) AuthFirebase(c *gin.Context) {}
func (ctrl *Auth) RefreshToken(c *gin.Context) {}
func (ctrl *Auth) Logout(c *gin.Context)       {}
