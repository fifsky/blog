package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"app/cmd"
	"app/config"
	"app/pkg/aiagent"
	"app/pkg/litestream"
	"app/service/feishu"
	"app/service/remind"
	"app/store"

	"github.com/goapt/logger"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"

	"github.com/urfave/cli/v3"
)

func main() {
	conf := config.New()

	// 1. 先启动 Litestream（作为 Go library 嵌入），从 OSS 自动恢复 + 实时备份 SQLite
	dbPath := conf.DB.ExtractDBPath()
	ls := litestream.New(conf, dbPath)
	if err := ls.Start(context.Background()); err != nil {
		log.Fatalf("[litestream] start failed: %s", err)
	}

	// 2. 再打开应用层数据库连接（确保 litestream 已初始化）
	db := conf.DB.Connect()

	defer func() {
		// 先关闭应用层数据库连接
		if err := db.Close(); err != nil {
			log.Printf("[db] database close error: %s", err)
		}
		// 再关闭 litestream store（确保所有应用连接已释放）
		_ = ls.Stop()
	}()

	// logger
	logger.SetDefault(config.NewLogger(conf, "app.log"))

	// crontab setup
	s := store.New(db)

	// 创建飞书发送器和卡片处理器（仅保留跨服务共享的部分）
	sender := feishu.NewFeishuSender(conf.Feishu)
	remindCard := feishu.NewRemindCard(s, conf)
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

	// TODO 临时停用自动生成心情
	// ai := motto.NewOpenAIProvider(agent)
	// m := motto.New(s, ai)
	// go m.Start("0 7 * * *")

	app := &cli.Command{
		Name:  "blog",
		Usage: "fifsky blog",
		Commands: []*cli.Command{
			cmd.NewHttp(s, conf, agent),
			cmd.NewTmp(db, conf),
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
