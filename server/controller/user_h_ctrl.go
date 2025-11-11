package controller

import "github.com/gin-gonic/gin"

type User struct{}

func NewUserController() *User {
	return &User{}
}

func (ctrl *User) GetProfile(c *gin.Context)     {}
func (ctrl *User) UpdateProfile(c *gin.Context)  {}
func (ctrl *User) ChangePassword(c *gin.Context) {}
