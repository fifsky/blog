package middleware

import (
	"github.com/google/wire"
)

type Middleware struct {
	AccessLog
	Recover
	Cors
	Limiter
	RemindAuth
	AuthLogin
}

var ProviderSet = wire.NewSet(NewCors, NewLimiter, NewRecover, NewRemindAuth, NewAccessLog, NewAuthLogin, wire.Struct(new(Middleware), "*"))
