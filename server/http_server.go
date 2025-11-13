package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"scs-auth-service/config"
	"scs-auth-service/database"
	"scs-auth-service/helpers"
	"scs-auth-service/o11y/clog"
	"scs-auth-service/server/controller"
	"scs-auth-service/server/middleware"
	"scs-auth-service/server/repository"
	"scs-auth-service/server/router"
	"scs-auth-service/server/service"
)

const (
	defaultHTTPHost          = "localhost"
	defaultHTTPPort          = 8080
	defaultReadHeaderTimeout = time.Second * 5
	defaultReadTimeout       = time.Second * 10
	defaultWriteTimeout      = time.Second * 10
	defaultIdleTimeout       = time.Second * 60
	defaultShutdownTimeout   = time.Second * 30
)

type HTTPServer struct {
	*http.Server

	Config      *config.Config
	PostgresDBs map[string]*gorm.DB
}

// NewHTTPServer HTTP APIs for Shadow Cat auth service.
//
//	@title			Shadow Cat Auth HTTP APIs
//	@version		1.0
//	@description	Backend HTTP APIs for Shadow Cat auth service.
//	@scheme			http https
func NewHTTPServer(
	cfg *config.Config,
	pgDBs map[string]*gorm.DB,
) *HTTPServer {
	authRepo := repository.NewAuthRepository(cfg, pgDBs[database.AuthDB])

	clientCredMW := middleware.NewClientCredentialMiddleware(cfg)
	authMW := middleware.NewAuthMiddleware(cfg, authRepo)

	authSrv := service.NewAuthService(cfg, authRepo)

	authCtrl := controller.NewAuthController(cfg, authSrv)
	userCtrl := controller.NewUserController()

	r := gin.Default()
	r.Use(helpers.StandardCORS())

	router.LoadRoutes(r, cfg,
		clientCredMW, authMW,
		authCtrl, userCtrl)

	host := helpers.OrDefault(cfg.HTTP.Host, defaultHTTPHost)
	port := helpers.OrDefault(cfg.HTTP.Port, defaultHTTPPort)
	return &HTTPServer{
		Server: &http.Server{
			Addr:              fmt.Sprintf("%s:%d", host, port),
			Handler:           r,
			ReadHeaderTimeout: helpers.OrDefault(cfg.HTTP.ReadHeaderTimeout, defaultReadHeaderTimeout),
			ReadTimeout:       helpers.OrDefault(cfg.HTTP.ReadTimeout, defaultReadTimeout),
			WriteTimeout:      helpers.OrDefault(cfg.HTTP.WriteTimeout, defaultWriteTimeout),
			IdleTimeout:       helpers.OrDefault(cfg.HTTP.IdleTimeout, defaultIdleTimeout),
		},
		Config: cfg,
	}
}

func (s *HTTPServer) Run() {
	clog.Log().Info("HTTP server started", zap.String("address", s.Addr))
	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		clog.Log().Fatal("Failed to start HTTP server", zap.Error(err))
	}
}

func (s *HTTPServer) Stop() {
	shutdownTimeout := helpers.OrDefault(s.Config.HTTP.ShutdownTimeout, defaultShutdownTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		clog.Log().Error("HTTP server forced to shutdown", zap.Error(err))
	} else {
		clog.Log().Info("HTTP server shut down gracefully")
	}
}
