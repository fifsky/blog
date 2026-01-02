package sloghttp

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

type RoundTripper struct {
	next   http.RoundTripper
	logger *slog.Logger
	config Config
}

func (rt *RoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	start := time.Now()

	requestID := r.Header.Get(RequestIDHeaderKey)
	if rt.config.WithRequestID {
		if requestID == "" {
			requestID = uuid.New().String()
			r.Header.Set(RequestIDHeaderKey, requestID)
		}
		r = r.WithContext(context.WithValue(r.Context(), requestIDCtxKey, requestID))
	}

	// dump request body
	br := newBodyReader(r.Body, RequestBodyMaxSize, rt.config.WithRequestBody)
	r.Body = br

	// dump response body
	// bw := newBodyWriter(w, ResponseBodyMaxSize, config.WithResponseBody)
	res, err := rt.next.RoundTrip(r)
	bw := newResponse(res, RequestBodyMaxSize, rt.config.WithResponseBody)

	// Make sure we create a map only once per request (in case we have multiple middleware instances)
	if v := r.Context().Value(customAttributesCtxKey); v == nil {
		r = r.WithContext(context.WithValue(r.Context(), customAttributesCtxKey, &sync.Map{}))
	}

	defer log(rt.logger, rt.config, r, bw, br, start)

	return res, err
}

// NewRoundTripper returns a `http.RoundTripper` that logs requests using slog.
func NewRoundTripper(logger *slog.Logger, next http.RoundTripper, config Config) *RoundTripper {
	return &RoundTripper{
		next:   next,
		logger: logger,
		config: config,
	}
}
