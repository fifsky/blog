package main

import (
	"context"
	"log"
	"os"

	"app/cmd"
	"app/config"
	"app/pkg/wechat"
	"app/remind"
	"app/store"

	"github.com/urfave/cli/v3"
)

func main() {
	conf := config.New()
	db, clean := config.NewBlogDB(conf)
	defer clean()

	robot := wechat.NewRobot(conf.Common.RobotToken)
	// crontab setup
	s := store.New(db)
	r := remind.New(s, conf, robot)
	go r.Start()

	app := &cli.Command{
		Name:  "blog",
		Usage: "fifsky blog",
		Commands: []*cli.Command{
			cmd.NewHttp(db, conf, robot),
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
