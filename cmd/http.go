package cmd

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"app/config"
	"app/pkg/wechat"
	"app/server"
	"app/server/router"
	"app/service/admin"
	"app/service/openapi"
	"app/store"

	"github.com/urfave/cli/v3"
)

func NewHttp(db *sql.DB, conf *config.Config, robot *wechat.Robot, httpClient *http.Client) *cli.Command {
	return &cli.Command{
		Name:  "http",
		Usage: "http command eg: ./app http --addr=:8080",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "addr",
				Usage: "http listen ip:port",
			},
		},
		Action: func(ctx context.Context, cli *cli.Command) error {
			if !cli.IsSet("addr") {
				_ = cli.Set("addr", ":8080")
			}
			log.Println("[Env] Run profile:" + conf.Env)
			s := store.New(db)
			apiService := openapi.New(s, conf, robot, httpClient)
			adminService := admin.New(s, conf)
			route := router.New(apiService, adminService, conf, s)
			return server.New(
				server.Handler(route.Handler()),
				server.Address(cli.String("addr")),
			).Start(ctx)
		},
	}
}
