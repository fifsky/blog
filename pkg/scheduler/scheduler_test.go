package scheduler

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

func TestSchedulerCreation(t *testing.T) {
	s := New()
	if s == nil {
		t.Fatal("New() returned nil")
	}
}

func TestSchedulerWithTimezone(t *testing.T) {
	s := New(WithTimezone("America/New_York"))
	if s == nil {
		t.Fatal("New() with timezone returned nil")
	}
}

func TestJobRegistration(t *testing.T) {
	s := New()

	job := &Job{
		Name:     "test-registration",
		Schedule: "0 * * * *",
		Handler:  func(_ context.Context) error { return nil },
	}

	if err := s.Register(job); err != nil {
		t.Fatalf("failed to register valid job: %v", err)
	}

	// Registering duplicate name should fail
	if err := s.Register(job); err == nil {
		t.Error("expected error when registering duplicate job name")
	}
}

func TestSchedulerStartStop(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := New()

		var runCount atomic.Int32
		job := &Job{
			Name:     "test-start-stop",
			Schedule: "* * * * * *", // Every second (6-field format)
			Handler: func(_ context.Context) error {
				runCount.Add(1)
				return nil
			},
		}

		if err := s.Register(job); err != nil {
			t.Fatalf("failed to register job: %v", err)
		}

		// Start scheduler
		if err := s.Start(); err != nil {
			t.Fatalf("failed to start scheduler: %v", err)
		}

		// 在假时钟环境中等待 2.5 秒，synctest 会自动推进时间
		time.Sleep(2500 * time.Millisecond)
		synctest.Wait()

		// Stop scheduler
		if err := s.Stop(t.Context()); err != nil {
			t.Fatalf("failed to stop scheduler: %v", err)
		}

		count := runCount.Load()
		// Should have run at least twice (at 0s and 1s, maybe 2s)
		if count < 2 {
			t.Errorf("expected job to run at least 2 times, ran %d times", count)
		}

		// 验证停止后不再执行（假时钟推进但调度器 goroutine 已退出）
		finalCount := count
		time.Sleep(1500 * time.Millisecond)
		synctest.Wait()
		if runCount.Load() != finalCount {
			t.Error("scheduler did not stop - job continued running")
		}
	})
}

func TestSchedulerWithMiddleware(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var executionLog []string
		var logMu sync.Mutex

		logger := &testLogger{
			onInfo: func(msg string, _ ...any) {
				logMu.Lock()
				executionLog = append(executionLog, fmt.Sprintf("INFO: %s", msg))
				logMu.Unlock()
			},
			onError: func(msg string, _ ...any) {
				logMu.Lock()
				executionLog = append(executionLog, fmt.Sprintf("ERROR: %s", msg))
				logMu.Unlock()
			},
		}

		s := New(WithMiddleware(
			Recovery(func(jobName string, r any) {
				logMu.Lock()
				executionLog = append(executionLog, fmt.Sprintf("PANIC: %s - %v", jobName, r))
				logMu.Unlock()
			}),
			Logging(logger),
		))

		job := &Job{
			Name:     "test-middleware",
			Schedule: "* * * * * *", // Every second
			Handler: func(_ context.Context) error {
				return nil
			},
		}

		if err := s.Register(job); err != nil {
			t.Fatalf("failed to register job: %v", err)
		}

		if err := s.Start(); err != nil {
			t.Fatalf("failed to start: %v", err)
		}

		// 在假时钟环境中等待 1.5 秒，synctest 会自动推进时间
		time.Sleep(1500 * time.Millisecond)
		synctest.Wait()

		if err := s.Stop(t.Context()); err != nil {
			t.Fatalf("failed to stop: %v", err)
		}

		logMu.Lock()
		defer logMu.Unlock()

		// Should have at least one start and one completion log
		hasStart := false
		hasCompletion := false
		for _, log := range executionLog {
			if strings.Contains(log, "Job started") {
				hasStart = true
			}
			if strings.Contains(log, "Job completed") {
				hasCompletion = true
			}
		}

		if !hasStart {
			t.Error("expected job start log")
		}
		if !hasCompletion {
			t.Error("expected job completion log")
		}
	})
}
