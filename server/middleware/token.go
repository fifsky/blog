package middleware

import (
	"net/http"
	"strings"

	"app/pkg/errors"
	"app/server/response"
)

type Token = func(next http.Handler) http.Handler

func NewToken(token string) Token {
	expected := strings.TrimSpace(token)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if expected == "" {
				response.Fail(w, errors.Unauthorized("TOKEN_NOT_CONFIGURED", "token 未配置"))
				return
			}

			actual := strings.TrimSpace(readBearerToken(r))
			if actual == "" {
				response.Fail(w, errors.Unauthorized("TOKEN_MISSING", "缺少 token"))
				return
			}

			if actual != expected {
				response.Fail(w, errors.Unauthorized("TOKEN_INVALID", "token 无效"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func readBearerToken(r *http.Request) string {
	v := strings.TrimSpace(r.Header.Get("Authorization"))
	if v == "" {
		return ""
	}

	if token, ok := strings.CutPrefix(v, "Bearer "); ok {
		return strings.TrimSpace(token)
	}
	if token, ok := strings.CutPrefix(v, "bearer "); ok {
		return strings.TrimSpace(token)
	}
	return ""
}
