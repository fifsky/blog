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
		name         string
		path         string
		method       string
		expectStatus int
	}{
		{
			name:         "未匹配的路径返回404",
			path:         "/blog/not-exist",
			method:       http.MethodGet,
			expectStatus: http.StatusNotFound,
		},
		{
			name:         "OPTIONS预检请求返回200",
			path:         "/blog/",
			method:       http.MethodOptions,
			expectStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &config.Config{Env: "dev"}
			accessLogger := slog.Default()
			r := router.New(&openapi.Service{}, &adminsvc.Service{}, conf, &store.Store{}, accessLogger)
			handler := r.Handler()

			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectStatus, rr.Code)
		})
	}
}
