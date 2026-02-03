package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHeader(t *testing.T) {
	// 模拟一个简单的 Handler，用于验证 Context 中是否包含 Headers
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers, ok := RequestHeaders(r.Context())
		assert.True(t, ok)
		assert.Equal(t, "test-value", headers.Get("X-Test-Header"))
		w.WriteHeader(http.StatusOK)
	})

	handler := NewHeader(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Test-Header", "test-value")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestClientIPFromContext(t *testing.T) {
	tests := []struct {
		name     string
		headers  map[string]string
		expected string
	}{
		{
			name: "优先使用 X-Real-IP",
			headers: map[string]string{
				"X-Real-IP":       "1.2.3.4",
				"X-Forwarded-For": "5.6.7.8, 9.10.11.12",
			},
			expected: "1.2.3.4",
		},
		{
			name: "没有 X-Real-IP 使用 X-Forwarded-For",
			headers: map[string]string{
				"X-Forwarded-For": "5.6.7.8, 9.10.11.12",
			},
			expected: "5.6.7.8",
		},
		{
			name:     "没有相关 Header",
			headers:  map[string]string{},
			expected: "",
		},
		{
			name: "X-Real-IP 为空字符串",
			headers: map[string]string{
				"X-Real-IP":       " ",
				"X-Forwarded-For": "5.6.7.8",
			},
			expected: "5.6.7.8",
		},
		{
			name: "X-Forwarded-For 包含空格",
			headers: map[string]string{
				"X-Forwarded-For": "  5.6.7.8  , 9.10.11.12",
			},
			expected: "5.6.7.8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := http.Header{}
			for k, v := range tt.headers {
				header.Set(k, v)
			}

			// 构造包含 Header 的 Context
			ctx := context.WithValue(context.Background(), headerContextKey{}, header)

			got := ClientIPFromContext(ctx)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestRequestHeaders_NoHeaders(t *testing.T) {
	ctx := context.Background()
	_, ok := RequestHeaders(ctx)
	assert.False(t, ok)
}
