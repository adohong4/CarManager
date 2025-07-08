package helpers

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(UserName string) (string, error) {
	expiration := time.Now().Add(24 * time.Hour)

	claims := &jwt.StandardClaims{
		ExpiresAt: expiration.Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   UserName,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("some_value"))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
