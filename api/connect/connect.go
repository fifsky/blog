package connect

import (
	"time"

	"github.com/goapt/logger"
	"github.com/google/wire"
	"github.com/ilibs/gosql/v2"
	"golang.org/x/time/rate"

	"app/config"
)

func Connect(c *config.Config) {
	// db connection
	_ = gosql.Connect(c.DB)
}

func NewDB() *gosql.DB {
	return gosql.Use("default")
}

func NewRateLimiter() *rate.Limiter {
	return rate.NewLimiter(rate.Every(time.Second*1), 100000) // 根据项目配置，可以不开启
}

type AccessLogger logger.ILogger

func NewAccessLogger() AccessLogger {
	newLog := logger.NewLogger(func(c *logger.Config) {
		c.LogName = "access"
		c.LogMode = config.App.Log.LogMode
		c.LogPath = config.App.Log.LogPath
		c.LogLevel = config.App.Log.LogLevel
		c.LogMaxFiles = config.App.Log.LogMaxFiles
		c.LogSentryDSN = ""
		c.LogDetail = false
	})
	return newLog
}

var ProviderSet = wire.NewSet(NewDB, NewRateLimiter, NewAccessLogger)
