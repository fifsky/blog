package motto

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	aimotto "app/service/motto"
	"app/store"
	"app/store/model"

	"github.com/goapt/logger"
	"github.com/robfig/cron/v3"
)

// Motto 定时生成心情日志的后台任务，基于 cron 调度。
// 实现 runner.Task + runner.Stopper：Start 启动 cron 并阻塞至 ctx 取消，Stop 停止 cron 并等待运行中的任务完成。
type Motto struct {
	store *store.Store
	ai    aimotto.AIProvider
	spec  string
	cron  *cron.Cron
}

// New 创建心情生成任务，spec 为 cron 调度表达式（如 "0 7 * * *"）
func New(s *store.Store, ai aimotto.AIProvider, spec string) *Motto {
	return &Motto{
		store: s,
		ai:    ai,
		spec:  spec,
	}
}

// Name 返回任务名
func (m *Motto) Name() string { return "motto" }

// Start 启动 cron 调度并阻塞直至 ctx 取消。cron 的停止由 Stop 方法负责。
func (m *Motto) Start(ctx context.Context) error {
	c := cron.New()
	_, err := c.AddFunc(m.spec, func() {
		if err := m.generateDailyMotto(); err != nil {
			logger.Error("generate daily motto error", slog.String("err", err.Error()))
		} else {
			logger.Info("generate daily motto success")
		}
	})
	if err != nil {
		logger.Error("motto cron add func error", slog.String("err", err.Error()))
		return fmt.Errorf("motto cron add func: %w", err)
	}
	m.cron = c
	c.Start()

	// 阻塞等待 ctx 取消，cron 的优雅停止交由 Stop 处理
	<-ctx.Done()
	return nil
}

// Stop 停止 cron 调度并等待运行中的任务完成，由 Runner 在逆序停止阶段调用。
func (m *Motto) Stop() error {
	if m.cron != nil {
		stopCtx := m.cron.Stop()
		<-stopCtx.Done()
	}
	return nil
}

func (m *Motto) generateDailyMotto() error {
	logger.Info("start generate daily motto")
	dateStr := time.Now().Format("2006-01-02")

	content, err := m.ai.Generate(context.Background(), aimotto.Prompt, dateStr)
	if err != nil {
		return err
	}

	if content == "" {
		return fmt.Errorf("generate daily motto empty")
	}

	// 写入数据库
	md := &model.Mood{
		Content:   content,
		UserId:    3, // 固定位AI用户生成
		CreatedAt: time.Now(),
	}

	if _, err := m.store.CreateMood(context.Background(), md); err != nil {
		return err
	}

	return nil
}
