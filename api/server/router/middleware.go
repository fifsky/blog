package router

import "net/http"

type Middleware = func(next http.Handler) http.Handler

type Chain struct {
	mws []Middleware
}

func (m Chain) Handle(next http.Handler) http.Handler {
	for i := len(m.mws) - 1; i >= 0; i-- {
		next = m.mws[i](next)
	}
	return next
}

func (m Chain) HandlerFunc(next http.HandlerFunc) http.Handler {
	return m.Handle(next)
}

func (m Chain) Append(mws ...Middleware) Chain {
	return Chain{append(m.mws, mws...)}
}

func Use(mws ...Middleware) Chain {
	return Chain{mws}
}
