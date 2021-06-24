package main

import (
	"app/config"
	"app/connect"
	"app/pkg/remind"
	_ "github.com/go-sql-driver/mysql"
	"github.com/goapt/gee"
	"github.com/goapt/golib/robot"
	"github.com/goapt/golib/robot/ding"
	"github.com/goapt/logger"
)

func main() {
	connect.Connect(config.App)

	// logger init
	logger.Setting(func(c *logger.Config) {
		c.LogName = "app"
		c.LogMode = config.App.Log.LogMode
		c.LogPath = config.App.Log.LogPath
		c.LogLevel = config.App.Log.LogLevel
		c.LogMaxFiles = config.App.Log.LogMaxFiles
		c.LogSentryDSN = config.App.Log.LogSentryDSN
		c.LogSentryType = "go." + config.App.AppName
		c.LogDetail = config.App.Log.LogDetail
	})

	robot.Init(ding.NewRobot())
	robot.SetToken(config.App.Common.RobotToken)

	// crontab setup
	go remind.StartCron()

	// server setup
	cmds := Initialize()
	gee.NewCliServer().Run(cmds)
}
