package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"scs-auth-service/config"
	"scs-auth-service/helpers"
	"scs-auth-service/models/response"
	"scs-auth-service/o11y/clog"
	"scs-auth-service/server/repository"
)

type Auth struct {
	Config   *config.Config
	AuthRepo *repository.Auth
}

func NewAuthMiddleware(cfg *config.Config, authRepo *repository.Auth) *Auth {
	return &Auth{
		Config:   cfg,
		AuthRepo: authRepo,
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

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &response.Resp{
			ReturnCode:    helpers.EInvalidAccessToken,
			ReturnMessage: helpers.Message(helpers.EInvalidAccessToken),
		})
		return
	}
	userSession, err := m.AuthRepo.RUserSessionWUserID(c.Request.Context(), nil, userID)
	if err != nil || userSession == nil {
		clog.Log().Warn("AuthMiddleware.Handle - user session not found",
			zap.String("user_id", claims.UserID), zap.Error(err))
		c.AbortWithStatusJSON(http.StatusUnauthorized, &response.Resp{
			ReturnCode:    helpers.EInvalidAccessToken,
			ReturnMessage: helpers.Message(helpers.EInvalidAccessToken),
		})
		return
	} else {
		if userSession.Token != claims.Token {
			clog.Log().Warn("AuthMiddleware.Handle - user session token mismatch",
				zap.String("user_id", claims.UserID), zap.String("session_token", userSession.Token), zap.String("token_claim", claims.Token))
			c.AbortWithStatusJSON(http.StatusUnauthorized, &response.Resp{
				ReturnCode:    helpers.EOtherSessionActive,
				ReturnMessage: helpers.Message(helpers.EOtherSessionActive),
			})
			return
		}
		if userSession.Revoked {
			clog.Log().Warn("AuthMiddleware.Handle - user session revoked",
				zap.String("user_id", claims.UserID))
			c.AbortWithStatusJSON(http.StatusUnauthorized, &response.Resp{
				ReturnCode:    helpers.EExpiredAccessToken,
				ReturnMessage: helpers.Message(helpers.EExpiredAccessToken),
			})
			return
		}
	}

	c.Set("is_authenticated", true)
	c.Set("user_id", claims.UserID)
	c.Set("role", claims.Role)

	c.Next()
}
