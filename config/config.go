package config

import "time"

type Config struct {
	Production bool     `mapstructure:"-"`
	HTTP       HTTP     `mapstructure:"http"     validate:"required"`
	Database   Database `mapstructure:"database" validate:"required"`
	Secrets    Secrets  `mapstructure:"secrets"  validate:"required"`
}

type HTTP struct {
	Host              string        `mapstructure:"host"`
	Port              int           `mapstructure:"port"`
	ReadHeaderTimeout time.Duration `mapstructure:"read-header-timeout-duration"`
	ReadTimeout       time.Duration `mapstructure:"read-timeout-duration"`
	WriteTimeout      time.Duration `mapstructure:"write-timeout-duration"`
	IdleTimeout       time.Duration `mapstructure:"idle-timeout-duration"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdown-timeout-duration"`
}
