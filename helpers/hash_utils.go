package helpers

import (
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

func EncodeClientSignature(clientID, clientSecret string) string {
	hash := sha256.Sum256([]byte(clientID + ":" + clientSecret))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func HashUsingBcrypt(input string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	return string(hash), err
}

func VerifyBcryptHash(hashedValue, plainValue string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedValue), []byte(plainValue))
	return err == nil
}
