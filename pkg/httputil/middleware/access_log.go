package middleware

import (
	"log/slog"
	"net/http"

	"app/pkg/logger/sloghttp"
)

func AccessLog(logger *slog.Logger) Middleware {
	return func(rt http.RoundTripper) http.RoundTripper {
		return sloghttp.NewRoundTripper(logger, rt, sloghttp.Config{
			DefaultLevel:       slog.LevelInfo,
			ClientErrorLevel:   slog.LevelWarn,
			ServerErrorLevel:   slog.LevelError,
			WithRequestID:      true,
			WithUserAgent:      true,
			WithRequestHeader:  true,
			WithRequestBody:    true,
			WithResponseHeader: true,
			WithResponseBody:   true,
		})
	}
}
