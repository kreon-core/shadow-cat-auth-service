package helpers

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessClaims struct {
	jwt.RegisteredClaims

	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Token  string `json:"token"`
}

type RefreshClaims struct {
	jwt.RegisteredClaims

	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Token  string `json:"token"`
}

const (
	defaultAccessTokenExpiry  = time.Hour * 1
	defaultRefreshTokenExpiry = time.Hour * 24 * 7
)

func GenerateJWTAccessToken(userID, role, sessionToken string,
	issuer string, jwtSecretKey []byte, tokenExpiry time.Duration,
) (string, error) {
	claims := &AccessClaims{
		UserID: userID,
		Role:   role,
		Token:  sessionToken,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: issuer,
			ID:     uuid.NewString(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(
				OrDefault(tokenExpiry, defaultAccessTokenExpiry),
			)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Subject:  userID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}

func ParseJWTAccessToken(tokenString string, jwtSecretKey []byte) (*jwt.Token, *AccessClaims, error) {
	claims := &AccessClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(jwtToken *jwt.Token) (any, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecretKey, nil
	})
	return token, claims, err
}

func GenerateJWTRefreshToken(userID, role, sessionToken string,
	issuer string, jwtSecretKey []byte, tokenExpiry time.Duration,
) (string, error) {
	claims := &RefreshClaims{
		UserID: userID,
		Role:   role,
		Token:  sessionToken,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: issuer,
			ID:     sessionToken,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(
				OrDefault(tokenExpiry, defaultRefreshTokenExpiry),
			)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Subject:  userID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecretKey)
	return signedToken, err
}

func ParseJWTRefreshToken(tokenString string, jwtSecretKey []byte) (*jwt.Token, *RefreshClaims, error) {
	claims := &RefreshClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(jwtToken *jwt.Token) (any, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecretKey, nil
	})
	return token, claims, err
}
