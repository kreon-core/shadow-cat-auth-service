package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"scs-auth-service/config"
	"scs-auth-service/helpers"
	"scs-auth-service/models/request"
	"scs-auth-service/models/response"
	"scs-auth-service/o11y/clog"
	"scs-auth-service/server/service"
)

type Auth struct {
	Config  *config.Config
	AuthSrv *service.Auth
}

func NewAuthController(cfg *config.Config, authSrv *service.Auth) *Auth {
	return &Auth{
		Config:  cfg,
		AuthSrv: authSrv,
	}
}

func (ctrl *Auth) Register(c *gin.Context) {
	var req request.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		clog.Log().Warn("Auth.Register - invalid request",
			clog.SafeAny("request", &req), zap.Error(err))
		c.JSON(http.StatusBadRequest, &response.Resp{
			ReturnCode:    helpers.EInvalidRequest,
			ReturnMessage: helpers.Message(helpers.EInvalidRequest),
		})
		return
	}

	authData, eCode, err := ctrl.AuthSrv.Register(c.Request.Context(), &req)
	if err != nil {
		if eCode == helpers.EDatabaseError {
			clog.Log().Error("Auth.Register - service error",
				clog.SafeAny("request", &req), zap.Int("error_code", eCode), zap.Error(err))
		} else {
			clog.Log().Warn("Auth.Register - service error",
				clog.SafeAny("request", &req), zap.Int("error_code", eCode), zap.Error(err))
		}
		c.JSON(http.StatusBadRequest, &response.Resp{
			ReturnCode:    eCode,
			ReturnMessage: helpers.Message(eCode),
		})
		return
	}

	clog.Log().Info("Auth.Register - user registered successfully",
		clog.SafeAny("request", &req),
		zap.String("user_id", authData.UserID))
	c.JSON(http.StatusCreated, &response.Resp{
		ReturnCode:    helpers.SResourceCreatedSuccessfully,
		ReturnMessage: helpers.Message(helpers.SResourceCreatedSuccessfully),
		Data:          authData,
	})
}

func (ctrl *Auth) Login(c *gin.Context) {
	var req request.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		clog.Log().Warn("Auth.Login - invalid request",
			clog.SafeAny("request", &req), zap.Error(err))
		c.JSON(http.StatusBadRequest, &response.Resp{
			ReturnCode:    helpers.EInvalidRequest,
			ReturnMessage: helpers.Message(helpers.EInvalidRequest),
		})
		return
	}

	authData, eCode, err := ctrl.AuthSrv.Login(c.Request.Context(), &req)
	if err != nil {
		if eCode == helpers.EDatabaseError {
			clog.Log().Error("Auth.Login - service error",
				clog.SafeAny("request", &req), zap.Int("error_code", eCode), zap.Error(err))
		} else {
			clog.Log().Warn("Auth.Login - service error",
				clog.SafeAny("request", &req), zap.Int("error_code", eCode), zap.Error(err))
		}
		c.JSON(http.StatusBadRequest, &response.Resp{
			ReturnCode:    eCode,
			ReturnMessage: helpers.Message(eCode),
		})
		return
	}

	clog.Log().Info("Auth.Login - user logged in successfully",
		clog.SafeAny("request", &req),
		zap.String("user_id", authData.UserID))
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
		clog.Log().Warn("Auth.Logout - invalid request",
			clog.SafeAny("request", &req), zap.Error(err))
		c.JSON(http.StatusBadRequest, &response.Resp{
			ReturnCode:    helpers.EInvalidRequest,
			ReturnMessage: helpers.Message(helpers.EInvalidRequest),
		})
		return
	}

	eCode, err := ctrl.AuthSrv.Logout(c.Request.Context(), &req)
	if err != nil {
		if eCode == helpers.EDatabaseError {
			clog.Log().Error("Auth.Logout - service error",
				clog.SafeAny("request", &req), zap.Int("error_code", eCode), zap.Error(err))
		} else {
			clog.Log().Warn("Auth.Logout - service error",
				clog.SafeAny("request", &req), zap.Int("error_code", eCode), zap.Error(err))
		}
		c.JSON(http.StatusBadRequest, &response.Resp{
			ReturnCode:    eCode,
			ReturnMessage: helpers.Message(eCode),
		})
		return
	}

	clog.Log().Info("Auth.Logout - user logged out successfully",
		clog.SafeAny("request", &req))
	c.JSON(http.StatusOK, &response.Resp{
		ReturnCode:    helpers.Success,
		ReturnMessage: helpers.Message(helpers.Success),
	})
}
