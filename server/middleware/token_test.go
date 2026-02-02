package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewToken(t *testing.T) {
	tests := []struct {
		name           string
		configToken    string
		authHeader     string
		expectStatus   int
		expectBody     string
		expectNextCall bool
	}{
		{
			name:           "成功 - 有效的 token",
			configToken:    "valid-token-123",
			authHeader:     "Bearer valid-token-123",
			expectStatus:   http.StatusOK,
			expectBody:     `{"code":10000,"msg":"success"}`,
			expectNextCall: true,
		},
		{
			name:           "失败 - token 未配置",
			configToken:    "",
			authHeader:     "Bearer any-token",
			expectStatus:   http.StatusUnauthorized,
			expectBody:     "TOKEN_NOT_CONFIGURED",
			expectNextCall: false,
		},
		{
			name:           "失败 - 缺少 token",
			configToken:    "valid-token",
			authHeader:     "",
			expectStatus:   http.StatusUnauthorized,
			expectBody:     "TOKEN_MISSING",
			expectNextCall: false,
		},
		{
			name:           "失败 - 无效的 token",
			configToken:    "valid-token",
			authHeader:     "Bearer invalid-token",
			expectStatus:   http.StatusUnauthorized,
			expectBody:     "TOKEN_INVALID",
			expectNextCall: false,
		},
		{
			name:           "成功 - 小写 bearer 前缀",
			configToken:    "my-token",
			authHeader:     "bearer my-token",
			expectStatus:   http.StatusOK,
			expectBody:     "",
			expectNextCall: true,
		},
		{
			name:           "失败 - 大写 BEARER 前缀不支持",
			configToken:    "my-token",
			authHeader:     "BEARER my-token",
			expectStatus:   http.StatusUnauthorized,
			expectBody:     "TOKEN_MISSING",
			expectNextCall: false,
		},
		{
			name:           "成功 - 前后空格被正确处理",
			configToken:    "  my-token  ",
			authHeader:     "Bearer   my-token   ",
			expectStatus:   http.StatusOK,
			expectBody:     "",
			expectNextCall: true,
		},
		{
			name:           "失败 - 没有 Bearer 前缀",
			configToken:    "my-token",
			authHeader:     "my-token",
			expectStatus:   http.StatusUnauthorized,
			expectBody:     "TOKEN_MISSING",
			expectNextCall: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 使用闭包追踪 next handler 是否被调用
			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"code":10000,"msg":"success"}`))
			})

			m := NewToken(tt.configToken)
			handler := m(next)
			req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectStatus, rr.Code)
			if tt.expectBody != "" {
				assert.Contains(t, rr.Body.String(), tt.expectBody)
			}
			assert.Equal(t, tt.expectNextCall, nextCalled)
		})
	}
}
