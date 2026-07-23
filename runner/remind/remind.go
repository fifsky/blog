package remind

import (
	"context"
	"log/slog"
	"time"

	"app/pkg/remindutil"
	"app/service/feishu"
	"app/store"
	"app/store/model"

	"github.com/goapt/logger"
)

// Remind 定时提醒轮询任务，每分钟扫描到期提醒并通过飞书卡片发送。
// 实现 runner.Task 接口，ctx 取消后优雅退出。
type Remind struct {
	store *store.Store
	card  *feishu.RemindCard
}

// New 创建提醒轮询任务，card 为发送提醒使用的飞书卡片处理器
func New(s *store.Store, card *feishu.RemindCard) *Remind {
	return &Remind{
		store: s,
		card:  card,
	}
}

// Name 返回任务名
func (r *Remind) Name() string { return "remind" }

// Start 启动定时提醒轮询，每 60 秒扫描一次，ctx 取消后退出。
func (r *Remind) Start(ctx context.Context) error {
	t := time.NewTicker(60 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case t1 := <-t.C:
			r.run(t1)
		}
	}
}

func (r *Remind) buildMessage(content string, v *model.Remind) feishu.RemindMessage {
	return feishu.RemindMessage{
		Content: content,
		Time:    time.Now().Format("2006-01-02 15:04"),
		ID:      v.Id,
	}
}

func (r *Remind) message(content string, v *model.Remind) {
	msg := r.buildMessage(content, v)
	if err := r.card.Send(context.Background(), msg); err != nil {
		logger.Error("remind send message error", slog.String("err", err.Error()))
	}
}

func (r *Remind) changeNextTime(v *model.Remind) {
	nextTime := remindutil.NextTimeFromRule(time.Now(), v)

	_ = r.store.UpdateRemindStatus(context.Background(), v.Id, model.RemindStatusPending)
	_ = r.store.UpdateRemindNextTime(context.Background(), v.Id, nextTime)
}

func (r *Remind) run(t time.Time) {
	reminds, _ := r.store.RemindAll(context.Background())

	for _, v := range reminds {
		content := v.Content

		// 如果是等待确认的消息，则每天都需要提醒
		if v.Status == model.RemindStatusPending {
			// 未确认的消息每天都需要在相同的时间点提醒
			if t.Format("15:04") == v.NextTime.Format("15:04") {
				v2 := v
				r.message(content, &v2)
				r.changeNextTime(&v2)
			}
			continue
		}

		if !v.NextTime.IsZero() && !t.Before(v.NextTime) {
			v2 := v
			r.message(content, &v2)
			// 如果发出提醒，在用户没有点击确认收到之前，会不断提醒，因此需要更新下一次提醒时间为次日相同时间点
			r.changeNextTime(&v2)
		}
	}
}
