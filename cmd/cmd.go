package cmd

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"sync"

	"app/config"
	"app/pkg/agent"
	"app/pkg/litestream"
	"app/service/feishu"
	"app/service/remind"
	"app/store"

	"github.com/goapt/logger"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/urfave/cli/v3"
)

type Command struct {
	store *store.Store
	conf  *config.Config
	agent *agent.Agent
	wg    sync.WaitGroup
}

func NewCommand() *Command {
	return &Command{}
}

func (c *Command) Init(ctx context.Context) (func(), error) {
	c.conf = config.New()
	// 1. 先启动 Litestream（作为 Go library 嵌入），从 OSS 自动恢复 + 实时备份 SQLite
	dbPath := c.conf.DB.ExtractDBPath()
	ls := litestream.New(c.conf.Litestream, c.conf.Env, dbPath)
	if err := ls.Start(ctx); err != nil {
		return nil, fmt.Errorf("[litestream] start failed: %w", err)
	}

	// 2. 再打开应用层数据库连接（确保 litestream 已初始化）
	db := c.conf.DB.Connect()
	c.store = store.New(db)

	c.agent = agent.New(
		agent.WithConfigProvider(func(ctx context.Context) (openai.Client, string) {
			aiCfg := c.store.GetAIConfig(ctx)
			logger.Debug("ai config", slog.Any("config", aiCfg))
			client := openai.NewClient(
				option.WithAPIKey(aiCfg.Token),
				option.WithBaseURL(aiCfg.Endpoint),
			)
			return client, aiCfg.Model
		}),
		agent.WithMCP(c.conf.MCP),
	)

	return func() {
		// 等待后台 goroutine（remind、feishu bot）退出，避免使用已关闭的 DB
		c.wg.Wait()
		// 先关闭应用层数据库连接
		if err := db.Close(); err != nil {
			log.Printf("[db] database close error: %s", err)
		}
		// 再关闭 litestream store（确保所有应用连接已释放）
		if err := ls.Stop(); err != nil {
			log.Printf("[litestream] stop failed: %s", err)
		}
	}, nil
}

func (c *Command) runRemind(ctx context.Context) {
	// 卡片处理器内部自行创建飞书发送器，无需外部注入
	remindCard := feishu.NewRemindCard(c.store, c.conf.Feishu)
	linkCard := feishu.NewLinkCard(c.store, c.conf.Feishu)

	r := remind.New(c.store, remindCard)
	c.wg.Go(func() {
		r.Start(ctx)
	})

	// 创建卡片注册表（用于 bot 回调分发）
	registry := feishu.NewCardRegistry()
	registry.Register(remindCard)
	registry.Register(linkCard)
	// Feishu bot service
	if c.conf.Feishu.Appid != "" {
		feishuBot := feishu.NewBot(c.conf.Feishu, c.store, c.agent, registry)
		c.wg.Go(func() {
			feishuBot.Start(ctx)
		})
	}
}

func (c *Command) runMotto(ctx context.Context) {
	// TODO 临时停用自动生成心情
	// ai := motto.NewOpenAIProvider(agent)
	// m := motto.New(s, ai)
	// go m.Start("0 7 * * *")
}

func (c *Command) Run(ctx context.Context) error {
	app := &cli.Command{
		Name:  "blog",
		Usage: "fifsky blog",
		Commands: []*cli.Command{
			httpCommand(c),
			tmpCommand(c),
		},
	}
	return app.Run(ctx, os.Args)
}
