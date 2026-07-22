package remind

import (
	"context"
	"log/slog"
	"time"

	"app/config"
	"app/pkg/remindutil"
	"app/service/feishu"
	"app/store"
	"app/store/model"

	"github.com/goapt/logger"
)

type Remind struct {
	store *store.Store
	conf  *config.Config
	card  *feishu.RemindCard
}

func New(s *store.Store, conf *config.Config, card *feishu.RemindCard) *Remind {
	return &Remind{
		store: s,
		conf:  conf,
		card:  card,
	}
}

// Start 启动定时提醒轮询，ctx 取消后优雅退出
func (r *Remind) Start(ctx context.Context) {
	t := time.NewTicker(60 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case t1 := <-t.C:
			r.run(t1)
		}
	}
}

// NextTimeFromRule 根据提醒规则计算下一次触发时间，委托给 remindutil 包
func NextTimeFromRule(from time.Time, m *model.Remind) time.Time {
	return remindutil.NextTimeFromRule(from, m)
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
	nextTime := NextTimeFromRule(time.Now(), v)

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
