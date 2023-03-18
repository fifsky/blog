package main

import (
	"app/config"
	"app/connect"
	"app/pkg/remind"
	_ "github.com/go-sql-driver/mysql"
	"github.com/goapt/gee"
	"github.com/goapt/golib/robot"
	"github.com/goapt/golib/robot/wechat"
	"github.com/goapt/logger"
)

func main() {
	conf := config.New()
	connect.Connect(conf)

	// logger init
	logger.Setting(func(c *logger.Config) {
		c.LogName = "app"
		c.LogMode = conf.Log.LogMode
		c.LogPath = conf.Log.LogPath
		c.LogLevel = conf.Log.LogLevel
		c.LogMaxFiles = conf.Log.LogMaxFiles
		c.LogSentryDSN = conf.Log.LogSentryDSN
		c.LogSentryType = "go." + conf.AppName
		c.LogDetail = conf.Log.LogDetail
	})

	robot.Init(wechat.NewRobot())
	robot.SetToken(conf.Common.RobotToken)

	// crontab setup
	go remind.StartCron(conf)

	// server setup
	cmds := Initialize(conf)
	gee.NewCliServer().Run(cmds)
}
