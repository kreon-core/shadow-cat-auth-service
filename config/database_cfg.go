package config

import "time"

type Database struct {
	GormOptions GormOptions         `mapstructure:"gorm-options"`
	Postgres    map[string]Postgres `mapstructure:"postgres"`
}

type GormOptions struct {
	SkipDefaultTransaction    bool          `mapstructure:"skip_default_transaction"`
	DefaultTransactionTimeout time.Duration `mapstructure:"default_transaction_timeout"`
	DefaultContextTimeout     time.Duration `mapstructure:"default_context_timeout"`
}

type Postgres struct {
	Host     string `mapstructure:"host"     validate:"required"`
	Port     int    `mapstructure:"port"     validate:"required"`
	User     string `mapstructure:"user"     validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	DBName   string `mapstructure:"dbname"   validate:"required"`
	SSLMode  string `mapstructure:"sslmode"`
	TimeZone string `mapstructure:"timezone"`

	WithoutReturning bool `mapstructure:"without-returning"`
}
