package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"scs-auth-service/config"
	"scs-auth-service/helpers"
	"scs-auth-service/models/response"
)

type Auth struct {
	Config *config.Config
}

func NewAuthMiddleware(cfg *config.Config) *Auth {
	return &Auth{
		Config: cfg,
	}
}

func (m *Auth) Handle(c *gin.Context) {
	if value, exists := c.Get("is_authenticated"); exists {
		if bl, ok := value.(bool); ok && bl {
			c.Next()
			return
		}
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &response.Resp{
			ReturnCode:    helpers.EMissingAuthorizationHeader,
			ReturnMessage: helpers.Message(helpers.EMissingAuthorizationHeader),
		})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, claims, err := helpers.ParseJWTAccessToken(tokenString, []byte(m.Config.Secrets.JWTSecretKey))
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &response.Resp{
				ReturnCode:    helpers.EExpiredAccessToken,
				ReturnMessage: helpers.Message(helpers.EExpiredAccessToken),
			})
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &response.Resp{
				ReturnCode:    helpers.EInvalidAccessToken,
				ReturnMessage: helpers.Message(helpers.EInvalidAccessToken),
			})
		}
		return
	}
	if !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &response.Resp{
			ReturnCode:    helpers.EInvalidAccessToken,
			ReturnMessage: helpers.Message(helpers.EInvalidAccessToken),
		})
		return
	}

	c.Set("is_authenticated", true)
	c.Set("user_id", claims.UserID)
	c.Set("role", claims.Role)

	c.Next()
}
