package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"sc-auth-service/config"
	"sc-auth-service/helpers"
	"sc-auth-service/models/response"
)

type ClientCredential struct {
	Config *config.Config
}

func NewClientCredentialMiddleware(cfg *config.Config) *ClientCredential {
	return &ClientCredential{
		Config: cfg,
	}
}

func (m *ClientCredential) Handle(c *gin.Context) {
	if value, exists := c.Get("is_authenticated"); exists {
		if bl, ok := value.(bool); ok && bl {
			c.Next()
			return
		}
	}

	clientID := c.GetHeader("X-Client-ID")
	signature := c.GetHeader("X-Client-Signature")

	if helpers.IsBlankString(&clientID) || helpers.IsBlankString(&signature) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &response.Resp{
			ReturnCode:    helpers.ESignatureInvalid,
			ReturnMessage: helpers.Message(helpers.ESignatureInvalid),
		})
		return
	}

	clientName, ok := m.validateClient(clientID, strings.TrimPrefix(signature, "Client "))
	if !ok || helpers.IsBlankString(&clientName) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &response.Resp{
			ReturnCode:    helpers.ESignatureInvalid,
			ReturnMessage: helpers.Message(helpers.ESignatureInvalid),
		})
		return
	}

	c.Set("is_authenticated", true)
	c.Set("client", clientName)

	c.Next()
}

func (m *ClientCredential) validateClient(clientID, signature string) (string, bool) {
	for name, cred := range m.Config.Secrets.ClientCredentials {
		if cred != nil && cred.ClientID == clientID {
			if signature == helpers.EncodeClientSignature(cred.ClientID, cred.ClientSecret) {
				return name, true
			}
		}
	}
	return "", false
}
