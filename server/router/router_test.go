package router_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"app/config"
	"app/server/router"
	adminsvc "app/service/admin"
	"app/service/openapi"
	"app/store"

	"github.com/stretchr/testify/assert"
)

func TestRouter_Handler(t *testing.T) {
	tests := []struct {
		name string
		path string
		// origin 为空表示不携带 Origin 头
		origin       string
		method       string
		expectStatus int
		// expectAllowOrigin 非空时校验 Access-Control-Allow-Origin 响应头
		expectAllowOrigin string
	}{
		{
			name:         "未匹配的路径返回404",
			path:         "/blog/not-exist",
			method:       http.MethodGet,
			expectStatus: http.StatusNotFound,
		},
		{
			name:              "已注册路径的OPTIONS预检请求由CORS收口返回204",
			path:              "/blog/article/list",
			method:            http.MethodOptions,
			origin:            "https://fifsky.com",
			expectStatus:      http.StatusNoContent,
			expectAllowOrigin: "https://fifsky.com",
		},
		{
			name:              "未注册路径的OPTIONS预检请求同样由CORS收口返回204",
			path:              "/blog/not-exist",
			method:            http.MethodOptions,
			origin:            "https://fifsky.com",
			expectStatus:      http.StatusNoContent,
			expectAllowOrigin: "https://fifsky.com",
		},
		{
			name:              "鉴权路由的OPTIONS预检请求不经过鉴权由CORS收口",
			path:              "/blog/admin/article/list",
			method:            http.MethodOptions,
			origin:            "https://fifsky.com",
			expectStatus:      http.StatusNoContent,
			expectAllowOrigin: "https://fifsky.com",
		},
		{
			name:         "非白名单来源的请求被拒绝",
			path:         "/blog/article/list",
			method:       http.MethodPost,
			origin:       "https://evil.example.com",
			expectStatus: http.StatusForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &config.Config{Env: "dev"}
			accessLogger := slog.Default()
			r := router.New(&openapi.Service{}, &adminsvc.Service{}, conf, &store.Store{}, accessLogger)
			handler := r.Handler()

			req := httptest.NewRequest(tt.method, tt.path, nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
				req.Header.Set("Access-Control-Request-Method", http.MethodPost)
			}
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectStatus, rr.Code)
			if tt.expectAllowOrigin != "" {
				assert.Equal(t, tt.expectAllowOrigin, rr.Header().Get("Access-Control-Allow-Origin"))
			}
		})
	}
}
