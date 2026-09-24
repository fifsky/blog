package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCors(t *testing.T) {
	tests := []struct {
		name string
		// origin 为空表示不携带 Origin 头
		origin string
		// sameOriginHost 为 true 时把 Origin 设置为与请求 Host 同源
		sameOriginHost bool
		method         string
		expectStatus   int
		expectOrigin   string
		expectNext     bool
	}{
		{
			name:         "白名单来源的普通请求正常下发并携带CORS响应头",
			origin:       "https://fifsky.com",
			method:       http.MethodGet,
			expectStatus: http.StatusOK,
			expectOrigin: "https://fifsky.com",
			expectNext:   true,
		},
		{
			name:         "通配子域名来源允许放行",
			origin:       "https://windiness.fifsky.com",
			method:       http.MethodGet,
			expectStatus: http.StatusOK,
			expectOrigin: "https://windiness.fifsky.com",
			expectNext:   true,
		},
		{
			name:         "本地开发端口来源允许放行",
			origin:       "http://localhost:5173",
			method:       http.MethodGet,
			expectStatus: http.StatusOK,
			expectOrigin: "http://localhost:5173",
			expectNext:   true,
		},
		{
			name:         "预检请求由中间件返回204且不下发",
			origin:       "https://fifsky.com",
			method:       http.MethodOptions,
			expectStatus: http.StatusNoContent,
			expectOrigin: "https://fifsky.com",
			expectNext:   false,
		},
		{
			name:         "无Origin的OPTIONS返回204且不下发",
			method:       http.MethodOptions,
			expectStatus: http.StatusNoContent,
			expectNext:   false,
		},
		{
			name:           "同源Origin的OPTIONS返回204且不下发",
			sameOriginHost: true,
			method:         http.MethodOptions,
			expectStatus:   http.StatusNoContent,
			expectNext:     false,
		},
		{
			name:         "非白名单来源返回403",
			origin:       "https://evil.example.com",
			method:       http.MethodGet,
			expectStatus: http.StatusForbidden,
			expectNext:   false,
		},
		{
			name:         "无Origin的普通请求正常下发",
			method:       http.MethodGet,
			expectStatus: http.StatusOK,
			expectNext:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			})
			handler := NewCors()(next)

			req := httptest.NewRequest(tt.method, "/blog/article/list", nil)
			switch {
			case tt.sameOriginHost:
				req.Header.Set("Origin", "http://"+req.Host)
			case tt.origin != "":
				req.Header.Set("Origin", tt.origin)
				req.Header.Set("Access-Control-Request-Method", http.MethodPost)
			}

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectStatus, rec.Code)
			assert.Equal(t, tt.expectNext, nextCalled)
			assert.Equal(t, tt.expectOrigin, rec.Header().Get("Access-Control-Allow-Origin"))
		})
	}
}

func TestCors_PreflightHeaders(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := NewCors()(next)

	req := httptest.NewRequest(http.MethodOptions, "/blog/article/list", nil)
	req.Header.Set("Origin", "https://fifsky.com")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)
	req.Header.Set("Access-Control-Request-Headers", "Authorization,Content-Type")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "true", rec.Header().Get("Access-Control-Allow-Credentials"))
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Methods"), http.MethodPost)
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Headers"), "Authorization")
	assert.Equal(t, "43200", rec.Header().Get("Access-Control-Max-Age"))
	assert.Contains(t, rec.Header().Values("Vary"), "Origin")
}

func TestCors_AllowAllOrigins(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := NewCors(CorsOrigins([]string{"*"}))(next)

	req := httptest.NewRequest(http.MethodGet, "/blog/article/list", nil)
	req.Header.Set("Origin", "https://any.example.com")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
}
