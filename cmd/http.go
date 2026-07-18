package cmd

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"time"

	"app/config"
	"app/pkg/aiagent"
	"app/pkg/httpx"
	"app/server"
	"app/server/router"
	"app/service/admin"
	"app/service/openapi"
	"app/store"

	"github.com/goapt/logger/sloghttp"
	"github.com/urfave/cli/v3"
)

func NewHttp(s *store.Store, conf *config.Config, agent *aiagent.Agent) *cli.Command {
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

			// 日志和 HTTP 客户端在 http 命令内就近创建，避免从 main 一路透传
			accessLogger := config.NewLogger(conf, "access.log")
			httpClient := httpx.NewClient(
				httpx.WithTimeout(10*time.Second),
				httpx.WithMiddleware(func(rt http.RoundTripper) http.RoundTripper {
					return sloghttp.NewRoundTripper(accessLogger, rt, sloghttp.Config{
						Level:              slog.LevelInfo,
						WithUserAgent:      true,
						WithRequestBody:    true,
						WithRequestHeader:  true,
						WithResponseBody:   true,
						WithResponseHeader: true,
					})
				}),
			)

			apiService := openapi.New(s, conf, httpClient)
			adminService := admin.New(s, conf, agent)
			route := router.New(apiService, adminService, conf, s, accessLogger)
			return server.New(
				server.Handler(route.Handler()),
				server.Address(cli.String("addr")),
				// AI 对话使用 SSE 长连接，不设置全局写入截止时间。
				server.WriteTimeout(0),
			).Start(ctx)
		},
	}
}
