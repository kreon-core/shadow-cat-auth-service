package dto

import "sc-auth-service/models/enum"

type AuthRefreshData struct {
	TokenType    enum.TokenType `json:"token_type"`
	AccessToken  string         `json:"access_token"`
	ExpiresIn    int64          `json:"expires_in"`
	RefreshToken string         `json:"refresh_token"`
}
