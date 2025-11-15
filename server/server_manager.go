package server

import (
	"fmt"
	"sync"

	"gorm.io/gorm"

	"sc-auth-service/config"
	"sc-auth-service/database"
	"sc-auth-service/models/entity"
	"sc-auth-service/o11y/clog"
)

type Manager struct {
	Config      *config.Config
	PostgresDBs map[string]*gorm.DB
	HTTPServer  chan Server
}

func NewServerManager(cfg *config.Config) *Manager {
	return &Manager{
		Config:     cfg,
		HTTPServer: make(chan Server, 1),
	}
}

func (sm *Manager) LoadPostgresDBs() error {
	sm.PostgresDBs = make(map[string]*gorm.DB)
	for name, cfg := range sm.Config.Database.Postgres {
		db, err := database.CreatePostgresConnection(
			&sm.Config.Database.GormOptions,
			&cfg,
		)
		if err != nil {
			return fmt.Errorf("init_postgres_db_%s -> %w", name, err)
		}
		sm.PostgresDBs[name] = db
	}

	err := database.LoadModels(sm.PostgresDBs[database.AuthDB],
		&entity.User{}, &entity.AuthProvider{}, &entity.UserSession{},
	)
	if err != nil {
		return fmt.Errorf("load_models -> %w", err)
	}

	return nil
}

func (sm *Manager) StartServers() {
	var wg sync.WaitGroup

	httpServer := NewHTTPServer(sm.Config, sm.PostgresDBs)
	sm.HTTPServer <- httpServer
	wg.Go(httpServer.Run)

	wg.Wait()
}

func (sm *Manager) ShutdownServers() {
	var wg sync.WaitGroup

	httpSrv := <-sm.HTTPServer
	if httpSrv != nil {
		wg.Go(httpSrv.Stop)
	}

	wg.Wait()
	clog.Log().Info("All servers shut down successfully")
}
