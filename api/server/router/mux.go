package router

import (
	"net/http"
)

type ServeMux struct {
	*http.ServeMux
	middlewares []Middleware
}

func NewServeMux() *ServeMux {
	return &ServeMux{
		ServeMux: http.NewServeMux(),
	}
}

func (s *ServeMux) Use(middlewares ...Middleware) *ServeMux {
	return &ServeMux{
		ServeMux:    s.ServeMux,
		middlewares: append(s.middlewares, middlewares...),
	}
}

func (s *ServeMux) Handle(pattern string, handler http.Handler) {
	s.ServeMux.Handle(pattern, chain(handler, s.middlewares))
}

func (s *ServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.ServeMux.Handle(pattern, chain(http.HandlerFunc(handler), s.middlewares))
}
