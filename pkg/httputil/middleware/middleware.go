package middleware

import (
	"net/http"
	"slices"
)

// Middleware is HTTP Client transport middleware.
type Middleware func(http.RoundTripper) http.RoundTripper

// Chain returns a Middleware that specifies the chained handler for endpoint.
func Chain(rt http.RoundTripper, middlewares ...Middleware) http.RoundTripper {
	if rt == nil {
		rt = http.DefaultTransport
	}

	for _, mw := range slices.Backward(middlewares) {
		rt = mw(rt)
	}

	return rt
}

// RoundTripFunc is a holder function to make the process of
// creating middleware a bit easier without requiring the consumer to
// implement the RoundTripper interface.
type RoundTripFunc func(*http.Request) (*http.Response, error)

func (rt RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt(req)
}
