package openapi

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func signAccessToken(tokenSecret string, userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		Issuer:    fmt.Sprint(userID),
	})
	return token.SignedString([]byte(tokenSecret))
}
