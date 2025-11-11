package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"scs-auth-service/config"
	"scs-auth-service/helpers"
	"scs-auth-service/models/request"
	"scs-auth-service/models/response"
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

func (ctrl *Auth) Register(c *gin.Context) {
	var req request.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &response.Resp{
			ReturnCode:    helpers.EInvalidRequest,
			ReturnMessage: helpers.Message(helpers.EInvalidRequest),
		})
		return
	}

	authData, eCode, err := ctrl.AuthSrv.Register(c.Request.Context(), nil, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, &response.Resp{
			ReturnCode:    eCode,
			ReturnMessage: helpers.Message(eCode),
		})
		return
	}

	c.JSON(http.StatusOK, &response.Resp{
		ReturnCode:    helpers.Success,
		ReturnMessage: helpers.Message(helpers.Success),
		Data:          authData,
	})
}

func (ctrl *Auth) Login(c *gin.Context) {
	var req request.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &response.Resp{
			ReturnCode:    helpers.EInvalidRequest,
			ReturnMessage: helpers.Message(helpers.EInvalidRequest),
		})
		return
	}

	authData, eCode, err := ctrl.AuthSrv.Login(c.Request.Context(), nil, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, &response.Resp{
			ReturnCode:    eCode,
			ReturnMessage: helpers.Message(eCode),
		})
		return
	}

	c.JSON(http.StatusOK, &response.Resp{
		ReturnCode:    helpers.Success,
		ReturnMessage: helpers.Message(helpers.Success),
		Data:          authData,
	})
}

func (ctrl *Auth) AuthZalo(c *gin.Context)     {}
func (ctrl *Auth) AuthFirebase(c *gin.Context) {}
func (ctrl *Auth) RefreshToken(c *gin.Context) {}

func (ctrl *Auth) Logout(c *gin.Context) {
	var req request.LogoutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &response.Resp{
			ReturnCode:    helpers.EInvalidRequest,
			ReturnMessage: helpers.Message(helpers.EInvalidRequest),
		})
		return
	}

	eCode, err := ctrl.AuthSrv.Logout(c.Request.Context(), nil, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, &response.Resp{
			ReturnCode:    eCode,
			ReturnMessage: helpers.Message(eCode),
		})
		return
	}

	c.JSON(http.StatusOK, &response.Resp{
		ReturnCode:    helpers.Success,
		ReturnMessage: helpers.Message(helpers.Success),
	})
}
