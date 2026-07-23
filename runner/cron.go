package runner

import (
	"context"
	"time"

	"app/pkg/scheduler"
)

// CronTask 将 scheduler.Scheduler 适配为 runner.Task + runner.Stopper。
// 用于全局统一管理 cron 调度器的启动与优雅停止，多个 cron 任务共享同一调度器实例，
// 各任务通过 scheduler.Scheduler.Register 注册 Job，生命周期由此处集中控制。
type CronTask struct {
	sched *scheduler.Scheduler
}

// NewCronTask 创建调度器生命周期管理任务。
// sched 须在 Start 调用前完成所有 Job 注册（scheduler.Start 后不支持动态注册）。
func NewCronTask(sched *scheduler.Scheduler) *CronTask {
	return &CronTask{sched: sched}
}

// Name 返回任务名
func (t *CronTask) Name() string { return "cron-scheduler" }

// Start 启动调度器并阻塞直至 ctx 取消。调度器的停止由 Stop 方法负责。
func (t *CronTask) Start(ctx context.Context) error {
	if err := t.sched.Start(); err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}

// Stop 优雅停止调度器，等待运行中的任务完成（最长 30 秒）。
func (t *CronTask) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return t.sched.Stop(ctx)
}
