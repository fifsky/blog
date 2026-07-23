package cmd

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"time"

	"app/config"
	"app/pkg/httpx"
	"app/server"
	"app/server/router"
	"app/service/admin"
	"app/service/openapi"

	"github.com/goapt/logger/sloghttp"
	"github.com/urfave/cli/v3"
)

func httpCommand(c *Command) *cli.Command {
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
			clean, err := c.Init(ctx)
			if err != nil {
				return err
			}
			defer clean()

			log.Println("[Env] Run profile:" + c.conf.Env)

			// 启动后台服务（remind 定时提醒、feishu bot），与 http server 并发运行
			bg := c.runBackground(ctx)
			defer func() {
				// 显式停止需释放资源的 task，并等待所有后台退出，避免使用已关闭的 DB
				bg.Stop()
				bg.Wait()
			}()

			// 日志和 HTTP 客户端在 http 命令内就近创建，避免从 main 一路透传
			accessLogger := config.NewLogger(c.conf, "access.log")
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

			apiService := openapi.New(c.store, c.conf, httpClient)
			adminService := admin.New(c.store, c.conf, c.agent)
			route := router.New(apiService, adminService, c.conf, c.store, accessLogger)
			return server.New(
				server.Handler(route.Handler()),
				server.Address(cli.String("addr")),
				// AI 对话使用 SSE 长连接，不设置全局写入截止时间。
				server.WriteTimeout(0),
			).Start(ctx)
		},
	}
}
