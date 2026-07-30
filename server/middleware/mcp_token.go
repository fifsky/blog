package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/auth"
)

// NewMCPToken 创建符合 MCP 协议规范的 Bearer Token 鉴权中间件。
//
// 与 NewToken 的区别：使用 SDK 自带的 auth.RequireBearerToken，
// 401 响应符合 MCP 协议规范（标准客户端可据此走 OAuth 发现流程），
// 并将 TokenInfo 注入请求上下文，MCP 工具函数可通过 req.Extra.TokenInfo 读取。
//
// 静态 token 无过期时间，因此开启 AllowMissingExpiration，
// 否则中间件会拒绝 Expiration 为零值的 token。
func NewMCPToken(token string) Token {
	expected := strings.TrimSpace(token)
	verifier := func(ctx context.Context, actual string, _ *http.Request) (*auth.TokenInfo, error) {
		if actual != expected {
			return nil, fmt.Errorf("%w: invalid token", auth.ErrInvalidToken)
		}
		return &auth.TokenInfo{}, nil
	}
	return auth.RequireBearerToken(verifier, &auth.RequireBearerTokenOptions{
		AllowMissingExpiration: true,
	})
}
