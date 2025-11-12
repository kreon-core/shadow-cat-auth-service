package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"scs-auth-service/config"
	"scs-auth-service/server/controller"
	"scs-auth-service/server/middleware"

	_ "scs-auth-service/docs/swagger" // for swagger docs
)

func LoadRoutes(
	r *gin.Engine,
	cfg *config.Config,
	clientCredMW *middleware.ClientCredential,
	authMW *middleware.Auth,
	authCtrl *controller.Auth,
	userCtrl *controller.User,
) {
	r.Use(middleware.CatchGlobalHTTPError)

	rMeta := r.Group("")
	rMeta.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	if !cfg.Production {
		rMeta.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	rv1 := r.Group("/api/v1")

	authGroup := rv1.Group("/auth")
	authGroup.Use(clientCredMW.Handle)
	loadAuthRoutes(authGroup, authCtrl)

	userGroup := rv1.Group("/user/:id")
	userGroup.Use(authMW.Handle)
	loadUserRoutes(userGroup, authCtrl, userCtrl)
}

func loadAuthRoutes(
	rg *gin.RouterGroup,
	authCtrl *controller.Auth,
) {
	rg.POST("/register", authCtrl.Register)
	rg.POST("/login", authCtrl.Login)
	rg.POST("/zalo", authCtrl.AuthZalo)
	rg.POST("/firebase", authCtrl.AuthFirebase)
	rg.POST("/refresh", authCtrl.AuthRefresh)
}

func loadUserRoutes(
	rg *gin.RouterGroup,
	authCtrl *controller.Auth,
	userCtrl *controller.User,
) {
	rg.PUT("", userCtrl.UpdateProfile)
	rg.GET("", userCtrl.GetProfile)
	rg.PUT("/password", userCtrl.ChangePassword)
	rg.POST("/logout", authCtrl.Logout)
}
