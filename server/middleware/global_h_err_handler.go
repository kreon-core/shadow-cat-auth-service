package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sc-auth-service/helpers"
	"sc-auth-service/models/response"
)

func CatchGlobalHTTPError(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &response.Resp{
				ReturnCode:    helpers.UUnspecifiedError,
				ReturnMessage: helpers.Message(helpers.UUnspecifiedError),
			})
		}
	}()
	c.Next()
}
