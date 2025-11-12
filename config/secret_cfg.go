package config

import "time"

type Secrets struct {
	JWTIssuer          string                       `mapstructure:"jwt-issuer"           validate:"required"`
	JWTSecretKey       string                       `mapstructure:"jwt-secret-key"       validate:"required"`
	AccessTokenExpiry  time.Duration                `mapstructure:"access-token-expiry"`
	RefreshTokenExpiry time.Duration                `mapstructure:"refresh-token-expiry"`
	ClientCredentials  map[string]*ClientCredential `mapstructure:"client-credentials"`
	ZaloSecret         ZaloSecret                   `mapstructure:"zalo-secret"          validate:"required"`
}
type ClientCredential struct {
	ClientID     string `mapstructure:"client-id"     validate:"required"`
	ClientSecret string `mapstructure:"client-secret" validate:"required"`
}

type ZaloSecret struct {
	NrmSignSecret string `mapstructure:"nrm-sign-secret" validate:"required"`
}
