package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/goapt/gee"
	"github.com/goapt/golib/robot"
	"github.com/goapt/golib/robot/ding"
	"github.com/ilibs/gosql/v2"

	"app/cmd"
	"app/config"
	"app/pkg/remind"
)

func main() {
	_ = gosql.Connect(config.App.DB)
	robot.Init(ding.NewRobot())
	robot.SetToken(config.App.Common.RobotToken)

	// 定时提醒
	go remind.StartCron()

	// command server
	cliServ := gee.NewCliServer()
	cliServ.Run(cmd.Commands())
}
