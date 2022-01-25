package middleware

import (
	"time"

	"app/config"
	"github.com/gin-contrib/cors"
	"github.com/goapt/gee"
)

type Cors gee.HandlerFunc

func NewCors(conf *config.Config) Cors {
	return func(c *gee.Context) gee.Response {
		origins := []string{"http://fifsky.com", "http://www.fifsky.com", "https://fifsky.com", "https://www.fifsky.com"}

		if conf.Env == "local" {
			origins = []string{"*"}
		}

		cors.New(cors.Config{
			AllowOrigins: origins,
			AllowMethods: []string{"*"},
			AllowHeaders: []string{
				"Origin",
				"Content-Length",
				"Content-Type",
				"Access-Token",
				"Access-Control-Allow-Origin",
				"x-requested-with",
			},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		})(c.Context)

		return nil
	}
}
