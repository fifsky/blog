// Package feishu 将飞书机器人（service/feishu.Bot）的启动包装为 runner.Task。
// service/feishu 包含卡片处理器、发送器等被多处复用的能力，保留在 service/ 下不迁移，
// 本包仅负责后台 WebSocket 长连接的生命周期管理。
package feishu

import (
	"context"

	"app/pkg/agent"
	"app/service/feishu"
)

// BotTask 飞书机器人后台任务，包装 service/feishu.Bot 的启动。
// 通过 WebSocket 长连接监听飞书消息与卡片回调，ctx 取消后退出。
type BotTask struct {
	bot *feishu.Bot
}

// New 创建飞书机器人任务，复用 service/feishu.NewBot。
// registry 为卡片回调分发注册表，由调用方装配并注册各业务卡片后传入。
func New(conf feishu.Config, aiAgent *agent.Agent, registry *feishu.CardRegistry) *BotTask {
	return &BotTask{bot: feishu.NewBot(conf, aiAgent, registry)}
}

// Name 返回任务名
func (b *BotTask) Name() string { return "feishu_bot" }

// Start 启动飞书机器人 WebSocket 连接，阻塞直至 ctx 取消。
func (b *BotTask) Start(ctx context.Context) error {
	b.bot.Start(ctx)
	return nil
}

// Stop 关闭飞书机器人 WebSocket 连接并禁用自动重连，由 Runner 在逆序停止阶段调用。
func (b *BotTask) Stop() error {
	b.bot.Close()
	return nil
}
