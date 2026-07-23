package remind

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"app/pkg/remindutil"
	"app/pkg/scheduler"
	"app/service/feishu"
	"app/store"
	"app/store/model"

	"github.com/goapt/logger"
)

// Remind 定时提醒轮询任务，每分钟扫描到期提醒并通过飞书卡片发送。
// 通过 New 注册到共享的 scheduler.Scheduler，调度器的生命周期由 runner.CronTask 统一管理。
type Remind struct {
	store *store.Store
	card  *feishu.RemindCard
}

// New 创建提醒轮询任务并注册到共享调度器，每分钟执行一次扫描。
// 调度器的启动与停止由外部（runner.CronTask）统一控制，本任务不参与生命周期管理。
func New(sched *scheduler.Scheduler, s *store.Store, card *feishu.RemindCard) (*Remind, error) {
	r := &Remind{store: s, card: card}
	if err := sched.Register(&scheduler.Job{
		Name:     "remind",
		Schedule: "* * * * *", // 每分钟扫描一次到期提醒
		Handler:  r.handler,
	}); err != nil {
		return nil, fmt.Errorf("remind register job: %w", err)
	}
	return r, nil
}

// handler 是调度器触发的任务处理函数，扫描到期提醒并发送飞书消息。
func (r *Remind) handler(ctx context.Context) error {
	r.run(ctx, time.Now())
	return nil
}

func (r *Remind) message(ctx context.Context, content string, v *model.Remind) {
	if err := r.card.Send(ctx, feishu.RemindMessage{
		Content: content,
		Time:    time.Now().Format("2006-01-02 15:04"),
		ID:      v.Id,
	}); err != nil {
		logger.Error("remind send message error", slog.String("err", err.Error()))
	}
}

func (r *Remind) changeNextTime(ctx context.Context, v *model.Remind) {
	nextTime := remindutil.NextTimeFromRule(time.Now(), v)

	_ = r.store.UpdateRemindStatus(ctx, v.Id, model.RemindStatusPending)
	_ = r.store.UpdateRemindNextTime(ctx, v.Id, nextTime)
}

func (r *Remind) run(ctx context.Context, t time.Time) {
	reminds, _ := r.store.RemindAll(ctx)

	for _, v := range reminds {
		content := v.Content

		// 如果是等待确认的消息，则每天都需要提醒
		if v.Status == model.RemindStatusPending {
			// 未确认的消息每天都需要在相同的时间点提醒
			if t.Format("15:04") == v.NextTime.Format("15:04") {
				v2 := v
				r.message(ctx, content, &v2)
				r.changeNextTime(ctx, &v2)
			}
			continue
		}

		if !v.NextTime.IsZero() && !t.Before(v.NextTime) {
			v2 := v
			r.message(ctx, content, &v2)
			// 如果发出提醒，在用户没有点击确认收到之前，会不断提醒，因此需要更新下一次提醒时间为次日相同时间点
			r.changeNextTime(ctx, &v2)
		}
	}
}
