package cmd

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"app/config"
	"app/pkg/agent"
	"app/pkg/litestream"
	"app/runner"
	"app/runner/feishu"
	"app/runner/remind"
	"app/service/feishu"
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
		// 后台任务已由调用方通过 runner.Stop/Wait 退出（见 httpCommand 中的 defer）
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

// runBackground 装配并启动后台任务（提醒轮询、飞书机器人），返回 runner 供调用方优雅停止。
// 卡片处理器与注册表在此集中装配并注入相关 task，保持依赖解耦。
func (c *Command) runBackground(ctx context.Context) *runner.Runner {
	// 卡片处理器内部自行创建飞书发送器，无需外部注入
	remindCard := feishu.NewRemindCard(c.store, c.conf.Feishu)
	linkCard := feishu.NewLinkCard(c.store, c.conf.Feishu)

	// 卡片注册表（用于 bot 回调分发），注册提醒与友情链接卡片
	registry := feishu.NewCardRegistry()
	registry.Register(remindCard)
	registry.Register(linkCard)

	r := runner.New()
	// 提醒定时轮询
	r.Register(remind.New(c.store, remindCard))
	// 飞书机器人（仅配置了 Appid 时启动）
	if c.conf.Feishu.Appid != "" {
		r.Register(feishubot.New(c.conf.Feishu, c.agent, registry))
	}
	// motto 临时停用；如需启用：
	//   import (
	//     "app/runner/motto"
	//     aimotto "app/service/motto"
	//   )
	//   r.Register(motto.New(c.store, aimotto.NewOpenAIProvider(c.agent), "0 7 * * *"))
	_ = r.Start(ctx)
	return r
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
