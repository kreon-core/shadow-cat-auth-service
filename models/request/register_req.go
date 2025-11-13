package request

type RegisterReq struct {
	Username    string `json:"username"     binding:"omitempty,min=3,max=64"`
	Email       string `json:"email"        binding:"omitempty,email,max=128"`
	Password    string `json:"password"     binding:"required,min=6,max=512"  sensitive:"true"`
	DisplayName string `json:"display_name" binding:"omitempty,max=128"`
	AvatarURL   string `json:"avatar_url"   binding:"omitempty,url"`
}

func (r *RegisterReq) Valid() bool {
	return r.Username != "" || r.Email != ""
}
