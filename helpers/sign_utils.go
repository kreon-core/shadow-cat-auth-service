package helpers

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"sort"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func EncodeClientSignature(clientID, clientSecret string) string {
	hash := sha256.Sum256([]byte(clientID + ":" + clientSecret))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func HashUsingBcrypt(input string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func VerifyBcryptHash(hashedValue, plainValue string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedValue), []byte(plainValue))
	return err == nil
}

func ValidateZaloDataWithZaloSign(data map[string]string, secretKey string) bool {
	sign, ok := data["sign"]
	if !ok || IsBlankString(&sign) {
		return false
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		if k == "sign" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var builder strings.Builder
	for _, k := range keys {
		v := strings.TrimSpace(data[k])
		builder.WriteString(v)
	}

	raw := secretKey + builder.String()
	hash := md5.Sum([]byte(raw)) //nolint:gosec // follows Zalo's requirement
	expectedSign := hex.EncodeToString(hash[:])
	return expectedSign == sign
}
