package remind

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"app/config"
	"app/pkg/aesutil"
	"app/service/feishu"
	"app/store"
	"app/store/model"

	"github.com/goapt/logger"
	"github.com/robfig/cron/v3"
)

type Remind struct {
	store  *store.Store
	conf   *config.Config
	card   *feishu.RemindCard
	sender *feishu.FeishuSender
}

func New(s *store.Store, conf *config.Config, card *feishu.RemindCard, sender *feishu.FeishuSender) *Remind {
	return &Remind{
		store:  s,
		conf:   conf,
		card:   card,
		sender: sender,
	}
}

func (r *Remind) Start() {
	t := time.NewTicker(60 * time.Second)

	for t1 := range t.C {
		r.run(t1)
	}
}

func NextTimeFromRule(from time.Time, m *model.Remind) time.Time {
	if m.Cron == "" {
		return time.Time{}
	}

	location, _ := time.LoadLocation("Asia/Shanghai")
	if location == nil {
		location = time.Local
	}
	base := from.In(location)

	// If Cron looks like a specific date "2006-01-02 15:04:00"
	if len(m.Cron) >= 10 && m.Cron[4] == '-' {
		t, err := time.ParseInLocation(time.DateTime, m.Cron, location)
		if err == nil {
			return t
		}
	}

	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(m.Cron)
	if err != nil {
		return time.Time{}
	}
	return schedule.Next(base)
}

func (r *Remind) buildMessage(content string, v *model.Remind) feishu.RemindMessage {
	token, _ := aesutil.AesEncode(r.conf.Common.TokenSecret, strconv.Itoa(v.Id))

	return feishu.RemindMessage{
		Content: content,
		Time:    time.Now().Format("2006-01-02 15:04"),
		Token:   token,
	}
}

func (r *Remind) message(content string, v *model.Remind) {
	if r.sender == nil {
		return
	}

	msg := r.buildMessage(content, v)
	cardJSON := r.card.BuildCard(msg)

	if err := r.sender.Send(context.Background(), cardJSON); err != nil {
		logger.Error("remind send message error", slog.String("err", err.Error()))
	}
}

func (r *Remind) changeNextTime(v *model.Remind) {
	nextTime := NextTimeFromRule(time.Now(), v)

	_ = r.store.UpdateRemindStatus(context.Background(), v.Id, 2)
	_ = r.store.UpdateRemindNextTime(context.Background(), v.Id, nextTime)
}

func (r *Remind) run(t time.Time) {
	reminds, _ := r.store.RemindAll(context.Background())

	for _, v := range reminds {
		content := v.Content

		// 如果是等待确认的消息，则每天都需要提醒
		if v.Status == 2 {
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
