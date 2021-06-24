package testutil

import (
	"app/config"
	"github.com/goapt/logger"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	logger.Setting(func(c *logger.Config) {
		c.LogMode = "std"
	})

	config.App.Common.TokenSecret = "abcdabcdabcdabcd"
}
