package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"app/cmd"
	"app/config"
	"app/pkg/aiagent"
	"app/pkg/httpx"
	"app/service/feishu"
	"app/service/motto"
	"app/service/openapi"
	"app/service/remind"
	"app/store"

	"github.com/goapt/logger"
	"github.com/goapt/logger/sloghttp"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"

	"github.com/urfave/cli/v3"
)

func main() {
	conf := config.New()
	db := conf.DB.Connect()
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("[db] database close error: %s", err)
		}
	}()

	// logger
	logger.SetDefault(config.NewLogger(conf, "app.log"))

	// accessLogger 用于HTTP请求日志和路由访问日志
	accessLogger := config.NewLogger(conf, "access.log")

	// httpClient
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

	// crontab setup
	s := store.New(db)

	// 创建飞书发送器和卡片处理器
	sender := feishu.NewFeishuSender(conf.Feishu)
	remindSvc := openapi.NewRemind(s, conf)
	remindCard := feishu.NewRemindCard(remindSvc, s, conf)
	linkCard := feishu.NewLinkCard(s, conf)

	// 创建卡片注册表（用于 bot 回调分发）
	registry := feishu.NewCardRegistry()
	registry.Register(remindCard)
	registry.Register(linkCard)

	r := remind.New(s, conf, remindCard, sender)
	go r.Start()

	agent := aiagent.New(
		aiagent.WithConfigProvider(func(ctx context.Context) (openai.Client, string) {
			aiCfg := s.GetAIConfig(ctx)
			logger.Debug("ai config", slog.Any("config", aiCfg))
			client := openai.NewClient(
				option.WithAPIKey(aiCfg.Token),
				option.WithBaseURL(aiCfg.Endpoint),
			)
			return client, aiCfg.Model
		}),
		aiagent.WithMCP(conf.MCP),
	)

	// Feishu bot service
	if conf.Feishu.Appid != "" {
		feishuBot := feishu.NewBot(conf, s, agent, registry)
		go feishuBot.Start(context.Background())
	}

	ai := motto.NewOpenAIProvider(agent)
	m := motto.New(s, ai)
	go m.Start("0 7 * * *")

	app := &cli.Command{
		Name:  "blog",
		Usage: "fifsky blog",
		Commands: []*cli.Command{
			cmd.NewHttp(s, conf, httpClient, agent, accessLogger, linkCard, sender),
			cmd.NewTmp(db, conf),
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
