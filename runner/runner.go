// Package runner 提供轻量的后台任务（background task）生命周期管理框架。
//
// 设计目标：统一管理多个后台任务的启动与停止，替代散落在 cmd 中的
// goroutine + WaitGroup 样板代码。调度方式（ticker/cron/长连接）由各 Task
// 内部自行决定，框架只负责：
//  1. 注册 Task
//  2. 为每个 Task 启动独立 goroutine 运行 Start(ctx)
//  3. 优雅停止：按注册逆序调用实现了 Stopper 的 Task.Stop()，并等待全部退出
//
// 约定：Task.Start(ctx) 必须阻塞运行直至 ctx 取消或发生错误，不得提前返回。
package runner

import (
	"context"
	"log/slog"
	"slices"
	"sync"

	"github.com/goapt/logger"
)

// Task 后台任务接口。
// Start 必须阻塞运行直至 ctx 取消或出错，框架为每个 Task 启动独立 goroutine 调用之。
type Task interface {
	// Name 返回任务名，用于日志标识
	Name() string
	// Start 启动任务并阻塞运行，ctx 取消后应优雅退出并返回
	Start(ctx context.Context) error
}

// Stopper 可选接口，需要显式释放资源（如关闭连接、停止调度器）的 Task 实现之。
// Runner.Stop 会按注册逆序调用已注册 Task 的 Stop。
// 对于纯 ctx 驱动的 Task（未实现 Stopper），靠 ctx 取消 + Wait 退出即可。
type Stopper interface {
	Stop() error
}

// Runner 管理多个后台 Task 的注册、启动、停止与退出等待。
type Runner struct {
	tasks []Task
	wg    sync.WaitGroup
}

// New 创建 Runner
func New() *Runner {
	return &Runner{}
}

// Register 注册一个后台 Task，按调用顺序记录（停止时逆序处理）。
func (r *Runner) Register(t Task) {
	r.tasks = append(r.tasks, t)
}

// Start 为每个已注册 Task 启动独立 goroutine 运行其 Start(ctx)。
// 单个 Task 的 Start 返回错误仅记录日志，不阻断其他 Task。
func (r *Runner) Start(ctx context.Context) error {
	for _, t := range r.tasks {
		r.wg.Go(func() {
			logger.Info("runner task started", slog.String("task", t.Name()))
			if err := t.Start(ctx); err != nil && err != context.Canceled {
				logger.Error("runner task exited with error",
					slog.String("task", t.Name()), slog.String("err", err.Error()))
				return
			}
			logger.Info("runner task stopped", slog.String("task", t.Name()))
		})
	}
	return nil
}

// Stop 按注册逆序调用实现了 Stopper 接口的 Task.Stop，用于显式释放资源。
// 对于纯 ctx 驱动的 Task（未实现 Stopper），此方法无操作。
func (r *Runner) Stop() {
	for _, v := range slices.Backward(r.tasks) {
		if s, ok := v.(Stopper); ok {
			if err := s.Stop(); err != nil {
				logger.Error("runner task stop error",
					slog.String("task", v.Name()), slog.String("err", err.Error()))
			}
		}
	}
}

// Wait 阻塞等待所有已启动 Task 的 goroutine 退出。
// 通常在 ctx 取消后调用，确保所有后台任务退出后再关闭底层资源（如 DB）。
func (r *Runner) Wait() {
	r.wg.Wait()
}
