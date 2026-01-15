package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRecover_Panic(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})
	handler := NewRecover(next)

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.JSONEq(t, `{"code":"SYSTEM_ERROR","message":"系统错误","details":{"cause":"boom"}}`, rr.Body.String())
}

func TestNewRecover_PassThrough(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":10000,"msg":"success"}`))
	})
	handler := NewRecover(next)

	req := httptest.NewRequest(http.MethodGet, "/ok", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `{"code":10000,"msg":"success"}`, rr.Body.String())
}
