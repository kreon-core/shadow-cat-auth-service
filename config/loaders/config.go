package loaders

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"scs-auth-service/config"
	"scs-auth-service/helpers"
	"scs-auth-service/o11y/clog"
)

func LoadConfig() (*config.Config, error) {
	var cfg config.Config

	env := os.Getenv("ENV")
	if !helpers.IsBlankString(&env) {
		if slices.Contains([]string{"prod", "stag", "production", "staging"}, env) {
			cfg.Production = true
		}

		err := LoadEnvFile(fmt.Sprintf(".env.%s", env))
		if err != nil {
			clog.Log().Warn("Unable to load additional environment variables",
				zap.String("env", env), zap.Error(err))
		}
	}

	v := viper.New()
	v.AddConfigPath(".")
	v.AddConfigPath("config")
	v.SetConfigType("yaml")

	v.SetConfigName("application")
	err := v.MergeInConfig()
	if err == nil {
		clog.Log().Info("Loaded default configuration")
	}

	if !helpers.IsBlankString(&env) {
		v.SetConfigName(fmt.Sprintf("application-%s", env))
		err = v.MergeInConfig()
		if err == nil {
			clog.Log().Info("Loaded specific configuration", zap.String("env", env))
		} else {
			clog.Log().Warn("No specific configuration found", zap.String("env", env))
		}
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	err = v.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal_failed -> %w", err)
	}

	validate := validator.New()
	err = validate.Struct(cfg)
	if err != nil {
		return nil, fmt.Errorf("config_validation -> %w", err)
	}

	return &cfg, nil
}
