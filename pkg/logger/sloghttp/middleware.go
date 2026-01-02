package sloghttp

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

// NewMiddleware returns a `func(http.Handler) http.Handler` (middleware) that logs requests using slog.
func NewMiddleware(logger *slog.Logger, config Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			requestID := r.Header.Get(RequestIDHeaderKey)
			if config.WithRequestID {
				if requestID == "" {
					requestID = uuid.New().String()
					r.Header.Set(RequestIDHeaderKey, requestID)
				}
				r = r.WithContext(context.WithValue(r.Context(), requestIDCtxKey, requestID))
			}

			// dump request body
			br := newBodyReader(r.Body, RequestBodyMaxSize, config.WithRequestBody)
			r.Body = br

			// dump response body
			bw := newBodyWriter(w, ResponseBodyMaxSize, config.WithResponseBody)

			// Make sure we create a map only once per request (in case we have multiple middleware instances)
			if v := r.Context().Value(customAttributesCtxKey); v == nil {
				r = r.WithContext(context.WithValue(r.Context(), customAttributesCtxKey, &sync.Map{}))
			}

			defer log(logger, config, r, bw, br, start)

			next.ServeHTTP(bw, r)
		})
	}
}
