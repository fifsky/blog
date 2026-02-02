package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewToken_Success(t *testing.T) {
	// 测试 token 验证通过的情况
	token := "valid-token-123"
	m := NewToken(token)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":10000,"msg":"success"}`))
	})

	handler := m(next)
	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `{"code":10000,"msg":"success"}`, rr.Body.String())
}

func TestNewToken_NotConfigured(t *testing.T) {
	// 测试 token 未配置的情况（空字符串）
	m := NewToken("")

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("next handler should not be called when token is not configured")
	})

	handler := m(next)
	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "TOKEN_NOT_CONFIGURED")
	assert.Contains(t, rr.Body.String(), "token 未配置")
}

func TestNewToken_Missing(t *testing.T) {
	// 测试缺少 token 的情况
	m := NewToken("valid-token")

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("next handler should not be called when token is missing")
	})

	handler := m(next)
	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	// 不设置 Authorization header
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "TOKEN_MISSING")
	assert.Contains(t, rr.Body.String(), "缺少 token")
}

func TestNewToken_Invalid(t *testing.T) {
	// 测试 token 无效的情况
	m := NewToken("valid-token")

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("next handler should not be called when token is invalid")
	})

	handler := m(next)
	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "TOKEN_INVALID")
	assert.Contains(t, rr.Body.String(), "token 无效")
}

func TestNewToken_CaseInsensitiveBearer(t *testing.T) {
	// 测试大小写不敏感的 Bearer 前缀
	token := "my-token"
	m := NewToken(token)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := m(next)

	// 测试小写 bearer
	req1 := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req1.Header.Set("Authorization", "bearer "+token)
	rr1 := httptest.NewRecorder()
	handler.ServeHTTP(rr1, req1)
	assert.Equal(t, http.StatusOK, rr1.Code)

	// 测试大写 BEARER
	req2 := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req2.Header.Set("Authorization", "BEARER "+token)
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	assert.Equal(t, http.StatusUnauthorized, rr2.Code)
}

func TestNewToken_WhitespaceTrim(t *testing.T) {
	// 测试空白字符被正确处理
	token := "my-token"
	// 配置 token 带前后空格
	m := NewToken("  " + token + "  ")

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := m(next)

	// 请求 token 也带前后空格
	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.Header.Set("Authorization", "Bearer   "+token+"   ")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestNewToken_NoBearerPrefix(t *testing.T) {
	// 测试没有 Bearer 前缀的情况（应该返回空 token）
	token := "my-token"
	m := NewToken(token)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("next handler should not be called")
	})

	handler := m(next)

	// Authorization header 没有 Bearer 前缀
	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.Header.Set("Authorization", token)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "TOKEN_MISSING")
}
