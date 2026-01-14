package httputil

import (
	"log/slog"
	"net/http"

	"github.com/goapt/logger"
	"github.com/goapt/logger/sloghttp"
)

func Debug() Middleware {
	return func(rt http.RoundTripper) http.RoundTripper {
		logger := logger.New(&logger.Config{
			Mode: logger.ModeStd,
		})

		return sloghttp.NewRoundTripper(logger, rt, sloghttp.Config{
			DefaultLevel:       slog.LevelDebug,
			ClientErrorLevel:   slog.LevelDebug,
			ServerErrorLevel:   slog.LevelDebug,
			WithRequestID:      true,
			WithUserAgent:      true,
			WithRequestHeader:  true,
			WithRequestBody:    true,
			WithResponseHeader: true,
		})
	}
}
