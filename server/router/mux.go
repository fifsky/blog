package router

import (
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"runtime"
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
	slog.Info(fmt.Sprintf("[http] %-25s  --> %s", pattern, nameOfFunction(handler)))
	s.ServeMux.Handle(pattern, chain(handler, s.middlewares))
}

func (s *ServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.Handle(pattern, http.HandlerFunc(handler))
}

func nameOfFunction(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
