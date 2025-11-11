package helpers

import (
	"crypto/sha256"
	"encoding/base64"
)

func ToClientCredSignature(clientID, clientSecret string) string {
	hash := sha256.Sum256([]byte(clientID + ":" + clientSecret))
	return base64.StdEncoding.EncodeToString(hash[:])
}
