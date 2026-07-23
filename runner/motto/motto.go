package motto

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"app/pkg/scheduler"
	"app/service/motto"
	"app/store"
	"app/store/model"

	"github.com/goapt/logger"
)

// AIProvider 定义 AI 接口，方便测试
type AIProvider interface {
	Generate(ctx context.Context, prompt, content string) (string, error)
}

// Motto 定时生成心情日志的 cron 任务。
// 通过 New 注册到共享的 scheduler.Scheduler，调度器的生命周期由 runner.CronTask 统一管理。
type Motto struct {
	store *store.Store
	ai    AIProvider
}

// New 创建心情生成任务并注册到共享调度器，spec 为 cron 调度表达式（如 "0 7 * * *"）。
// 调度器的启动与停止由外部（runner.CronTask）统一控制，本任务不参与生命周期管理。
func New(sched *scheduler.Scheduler, s *store.Store, ai AIProvider, spec string) (*Motto, error) {
	m := &Motto{store: s, ai: ai}
	if err := sched.Register(&scheduler.Job{
		Name:     "motto",
		Schedule: spec,
		Handler:  m.handler,
	}); err != nil {
		return nil, fmt.Errorf("motto register job: %w", err)
	}
	return m, nil
}

// handler 是调度器触发的任务处理函数，生成每日心情并写入数据库。
func (m *Motto) handler(_ context.Context) error {
	if err := m.generateDailyMotto(); err != nil {
		logger.Error("generate daily motto error", slog.String("err", err.Error()))
	} else {
		logger.Info("generate daily motto success")
	}
	return nil
}

func (m *Motto) generateDailyMotto() error {
	logger.Info("start generate daily motto")
	dateStr := time.Now().Format("2006-01-02")

	content, err := m.ai.Generate(context.Background(), motto.Prompt, dateStr)
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
