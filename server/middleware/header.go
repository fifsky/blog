package middleware

import (
	"context"
	"net/http"
	"strings"
)

type headerContextKey struct{}

func NewHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 将请求头写入上下文
		ctx := context.WithValue(r.Context(), headerContextKey{}, r.Header.Clone())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestHeaders(ctx context.Context) (http.Header, bool) {
	value := ctx.Value(headerContextKey{})
	if value == nil {
		return nil, false
	}
	headers, ok := value.(http.Header)
	if !ok {
		return nil, false
	}
	return headers, true
}

func ClientIPFromContext(ctx context.Context) string {
	headers, ok := RequestHeaders(ctx)
	if !ok {
		return ""
	}
	ip := strings.TrimSpace(headers.Get("X-Real-IP"))
	if ip != "" {
		return ip
	}
	forwarded := strings.TrimSpace(headers.Get("X-Forwarded-For"))
	if forwarded == "" {
		return ""
	}
	parts := strings.Split(forwarded, ",")
	return strings.TrimSpace(parts[0])
}
