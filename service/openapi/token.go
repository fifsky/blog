package openapi

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// signAccessToken 生成 JWT access token，同时返回过期时间的 Unix 秒级时间戳
func signAccessToken(tokenSecret string, userID int) (string, int64, error) {
	expiry := time.Now().Add(24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiry),
		Issuer:    fmt.Sprint(userID),
	})
	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", 0, err
	}
	return tokenString, expiry.Unix(), nil
}
