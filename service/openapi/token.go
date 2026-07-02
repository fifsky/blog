package openapi

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// defaultTokenExpiry 默认 token 有效期（Web / 小程序等）
	defaultTokenExpiry = 24 * time.Hour
	// iosTokenExpiry iOS 客户端 token 有效期（30 天）
	iosTokenExpiry = 30 * 24 * time.Hour
)

// tokenExpiryForClient 根据客户端类型返回对应的 token 有效期
// iOS 客户端发放 30 天有效期，其他保持 24 小时
func tokenExpiryForClient(clientType string) time.Duration {
	if clientType == "ios" {
		return iosTokenExpiry
	}
	return defaultTokenExpiry
}

// signAccessToken 生成 JWT access token，同时返回过期时间的 Unix 秒级时间戳
// expiry 由调用方决定（如 iOS 端 30 天，Web/小程序 24 小时）
func signAccessToken(tokenSecret string, userID int, expiry time.Duration) (string, int64, error) {
	expireTime := time.Now().Add(expiry)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expireTime),
		Issuer:    fmt.Sprint(userID),
	})
	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", 0, err
	}
	return tokenString, expireTime.Unix(), nil
}
