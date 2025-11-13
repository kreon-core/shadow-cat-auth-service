package request

type LoginReq struct {
	Username string `json:"username" binding:"omitempty,min=3,max=64"`
	Email    string `json:"email"    binding:"omitempty,email,max=128"`
	Password string `json:"password" binding:"required"                sensitive:"true"`
}

func (r *LoginReq) Valid() bool {
	return r.Username != "" || r.Email != ""
}
