package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"app/cmd"
	"app/config"
	"app/pkg/bark"
	"app/pkg/wechat"
	"app/service/motto"
	"app/service/remind"
	"app/store"

	"github.com/goapt/httpx"
	"github.com/goapt/logger"

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
	httpClient := httpx.NewClient(httpx.WithMiddleware(httpx.AccessLog(logger.Default())))
	barkClient := bark.New(httpClient, conf.Common.NotifyUrl, conf.Common.NotifyToken)

	robot := wechat.NewRobot(conf.Common.RobotToken)
	// crontab setup
	s := store.New(db)
	r := remind.New(s, conf, robot, barkClient)
	go r.Start()

	ai := motto.NewDoubaoProvider(conf.Common.AIToken, conf.Common.AIModel)
	m := motto.New(s, conf, barkClient, ai)
	go m.Start("0 7,10,22 * * *")

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
