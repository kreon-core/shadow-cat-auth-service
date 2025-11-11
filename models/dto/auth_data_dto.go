package dto

import "time"

type AuthData struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`

	Role        string `json:"role"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	Status      string `json:"status"`

	Provider *AuthProviderData `json:"provider,omitempty"`

	PlayerID string `json:"player_id"`

	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type AuthProviderData struct {
	Name     string    `json:"name"`
	UID      string    `json:"uid"`
	LinkedAt time.Time `json:"linked_at"`
	Metadata any       `json:"metadata,omitempty"`
}
