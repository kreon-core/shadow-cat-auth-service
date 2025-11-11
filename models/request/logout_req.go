package request

type LogoutReq struct {
	UserID       string `json:"user_id"       binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}
