package cmd

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"app/config"
	"app/pkg/aiagent"
	"app/server"
	"app/server/router"
	"app/service/admin"
	"app/service/openapi"
	"app/store"

	"github.com/urfave/cli/v3"
)

func NewHttp(s *store.Store, conf *config.Config, httpClient *http.Client, agent *aiagent.Agent, accessLogger *slog.Logger) *cli.Command {
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

			apiService := openapi.New(s, conf, httpClient)
			adminService := admin.New(s, conf, agent)
			route := router.New(apiService, adminService, conf, s, accessLogger)
			return server.New(
				server.Handler(route.Handler()),
				server.Address(cli.String("addr")),
			).Start(ctx)
		},
	}
}
