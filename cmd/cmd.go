package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"

	"app/config"
	"app/pkg/agent"
	"app/pkg/litestream"
	"app/pkg/scheduler"
	"app/runner"
	feishubot "app/runner/feishu"
	"app/runner/remind"
	"app/service/feishu"
	"app/store"

	"github.com/goapt/logger"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/urfave/cli/v3"
)

type App struct {
	store *store.Store
	conf  *config.Config
	agent *agent.Agent
	db    *sql.DB
	ls    *litestream.Manager
}

func NewApp() *App {
	return &App{}
}

func (a *App) Init(ctx context.Context) error {
	a.conf = config.New()
	// 1. 先启动 Litestream（作为 Go library 嵌入），从 OSS 自动恢复 + 实时备份 SQLite
	dbPath := a.conf.DB.ExtractDBPath()
	a.ls = litestream.New(a.conf.Litestream, a.conf.Env, dbPath)
	if err := a.ls.Start(ctx); err != nil {
		return fmt.Errorf("[litestream] start failed: %w", err)
	}

	// 2. 再打开应用层数据库连接（确保 litestream 已初始化）
	a.db = a.conf.DB.Connect()
	a.store = store.New(a.db)

	a.agent = agent.New(
		agent.WithConfigProvider(func(ctx context.Context) (openai.Client, string) {
			aiCfg := a.store.GetAIConfig(ctx)
			logger.Debug("ai config", slog.Any("config", aiCfg))
			client := openai.NewClient(
				option.WithAPIKey(aiCfg.Token),
				option.WithBaseURL(aiCfg.Endpoint),
			)
			return client, aiCfg.Model
		}),
		agent.WithMCP(a.conf.MCP),
	)

	return nil
}

// Close 释放 Init 阶段申请的资源（数据库连接、litestream）。
// 后台任务需由调用方通过 runner.Stop/Wait 先行退出（见 httpCommand 中的 defer），
// 之后再调用本方法，避免使用已关闭的 DB。重复调用安全。
func (a *App) Close() {
	// 先关闭应用层数据库连接
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			log.Printf("[db] database close error: %s", err)
		}
		a.db = nil
	}
	// 再关闭 litestream store（确保所有应用连接已释放）
	if a.ls != nil {
		if err := a.ls.Stop(); err != nil {
			log.Printf("[litestream] stop failed: %s", err)
		}
		a.ls = nil
	}
}

// runBackground 装配并启动后台任务（提醒轮询、飞书机器人），返回 runner 供调用方优雅停止。
// 卡片处理器与注册表在此集中装配并注入相关 task，保持依赖解耦。
func (a *App) runBackground(ctx context.Context) *runner.Runner {
	// 卡片处理器内部自行创建飞书发送器，无需外部注入
	remindCard := feishu.NewRemindCard(a.store, a.conf.Feishu)
	linkCard := feishu.NewLinkCard(a.store, a.conf.Feishu)

	// 卡片注册表（用于 bot 回调分发），注册提醒与友情链接卡片
	registry := feishu.NewCardRegistry()
	registry.Register(remindCard)
	registry.Register(linkCard)

	r := runner.New()
	// 共享 cron 调度器，多个 cron 任务统一注册到此实例，由 CronTask 集中启停
	sched := scheduler.New(scheduler.WithTimezone("Asia/Shanghai"))
	// 提醒定时轮询（每分钟扫描到期提醒）
	if _, err := remind.New(sched, a.store, remindCard); err != nil {
		logger.Error("remind register error", slog.String("err", err.Error()))
	}
	// 飞书机器人（仅配置了 Appid 时启动）
	if a.conf.Feishu.Appid != "" {
		r.Register(feishubot.New(a.conf.Feishu, a.agent, registry))
	}
	// motto 临时停用；如需启用：
	//   import (
	//     "app/runner/motto"
	//     aimotto "app/service/motto"
	//   )
	//   if _, err := motto.New(sched, a.store, aimotto.NewOpenAIProvider(a.agent), "0 7 * * *"); err != nil {
	//       logger.Error("motto register error", slog.String("err", err.Error()))
	//   }
	// cron 调度器作为 runner.Task 统一启停（须在所有 Job 注册之后）
	r.Register(runner.NewCronTask(sched))
	_ = r.Start(ctx)
	return r
}

func (a *App) Run(ctx context.Context) error {
	app := &cli.Command{
		Name:  "blog",
		Usage: "fifsky blog",
		Commands: []*cli.Command{
			httpCommand(a),
			tmpCommand(a),
		},
	}
	return app.Run(ctx, os.Args)
}
