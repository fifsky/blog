package router

import (
	"net/http"
	"slices"
)

type Middleware = func(next http.Handler) http.Handler

// chain builds a http.Handler composed of an inline middleware stack and endpoint
// handler in the order they are passed.
func chain(endpoint http.Handler, middlewares []Middleware) http.Handler {
	for _, middleware := range slices.Backward(middlewares) {
		endpoint = middleware(endpoint)
	}
	return endpoint
}
