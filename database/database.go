package database

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"scs-auth-service/config"
	"scs-auth-service/helpers"
	"scs-auth-service/o11y/clog"
)

const (
	AuthDB = "auth"
)

const (
	defaultDefaultTransactionTimeout = 10 * time.Second
	defaultDefaultContextTimeout     = 5 * time.Second
)

func CreatePostgresConnection(gormOptions *config.GormOptions, cfg *config.Postgres) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName,
		cfg.Port, cfg.SSLMode, cfg.TimeZone,
	)
	pgCfg := postgres.Config{
		DSN:              dsn,
		WithoutReturning: helpers.OrDefault(cfg.WithoutReturning, false),
	}
	db, err := gorm.Open(postgres.New(pgCfg), &gorm.Config{
		SkipDefaultTransaction: helpers.OrDefault(gormOptions.SkipDefaultTransaction, false),
		DefaultTransactionTimeout: helpers.OrDefault(
			gormOptions.DefaultTransactionTimeout,
			defaultDefaultTransactionTimeout,
		),
		DefaultContextTimeout: helpers.OrDefault(gormOptions.DefaultContextTimeout, defaultDefaultContextTimeout),
	})
	if err != nil {
		return nil, fmt.Errorf("postgres_connection -> %w", err)
	}
	clog.Log().Info("Successfully connected to Postgres",
		zap.String("host", cfg.Host), zap.String("dbname", cfg.DBName))
	return db, nil
}

func LoadModels(db *gorm.DB, models ...any) error {
	err := db.AutoMigrate(models...)
	if err != nil {
		return fmt.Errorf("auto_migrate -> %w", err)
	}
	clog.Log().Info("Database models migrated successfully")
	return nil
}
