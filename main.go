package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"app/cmd"
	"app/config"
	"app/pkg/httputil"
	"github.com/goapt/logger"
	"app/pkg/wechat"
	"app/remind"
	"app/store"

	"github.com/urfave/cli/v3"
)

func main() {
	conf := config.New()
	db, clean := config.NewBlogDB(conf)
	defer clean()

	// logger
	logger.SetDefault(logger.New(&logger.Config{
		Mode:     logger.ModeFile,
		FileName: filepath.Join(conf.Common.StoragePath, "logs", "app.log"),
		MaxFiles: 3,
		Detail:   true,
	}))

	// httpClient
	httpClient := httputil.NewClient(httputil.WithMiddleware(httputil.AccessLog(logger.Default())))

	robot := wechat.NewRobot(conf.Common.RobotToken)
	// crontab setup
	s := store.New(db)
	r := remind.New(s, conf, robot)
	go r.Start()

	app := &cli.Command{
		Name:  "blog",
		Usage: "fifsky blog",
		Commands: []*cli.Command{
			cmd.NewHttp(db, conf, robot, httpClient),
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
