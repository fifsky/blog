package remind

import (
	"context"
	"log/slog"
	"net/url"
	"strconv"
	"time"

	"app/config"
	"app/pkg/aesutil"
	"app/pkg/messenger"
	"app/store"
	"app/store/model"

	"github.com/goapt/logger"
	"github.com/robfig/cron/v3"
)

type Remind struct {
	store     *store.Store
	conf      *config.Config
	messenger messenger.Sender
}

func New(s *store.Store, conf *config.Config, msgSender messenger.Sender) *Remind {
	return &Remind{
		store:     s,
		conf:      conf,
		messenger: msgSender,
	}
}

func (r *Remind) Start() {
	t := time.NewTicker(60 * time.Second)

	for t1 := range t.C {
		r.run(t1)
	}
}

func NextTimeFromRule(from time.Time, m *model.Remind) time.Time {
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
		return base
	}
	return schedule.Next(base)
}

func (r *Remind) buildMessage(title, content string, v *model.Remind) messenger.Message {
	token, _ := aesutil.AesEncode(r.conf.Common.TokenSecret, strconv.Itoa(v.Id))

	// For non-Feishu channels, still use URL
	changeURL := "https://api.fifsky.com/blog/remind/change?token=" + url.QueryEscape(token)
	delayURL := "https://api.fifsky.com/blog/remind/delay?token=" + url.QueryEscape(token)

	return messenger.Message{
		Title:   title,
		Content: content,
		Time:    time.Now().Format("2006-01-02 15:04"),
		Token:   token,
		Actions: []messenger.Action{
			{Title: "标记完成", URL: changeURL},
			{Title: "稍后提醒", URL: delayURL},
		},
	}
}

func (r *Remind) message(title, content string, v *model.Remind) {
	msg := r.buildMessage(title, content, v)

	if err := r.messenger.Send(context.Background(), msg); err != nil {
		logger.Default().Error("remind send message error", slog.String("err", err.Error()))
	}
}

func (r *Remind) changeNextTime(id int) {
	location, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(location)
	nextTime := now.AddDate(0, 0, 1)

	_ = r.store.UpdateRemindStatus(context.Background(), id, 2)
	_ = r.store.UpdateRemindNextTime(context.Background(), id, nextTime)
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
				r.message("🙋🏻‍再次提醒", content, &v2)
				r.changeNextTime(v.Id)
			}
			continue
		}

		if !v.NextTime.IsZero() && !t.Before(v.NextTime) {
			v2 := v
			r.message("⏰重要提醒⏰", content, &v2)
			// 如果发出提醒，在用户没有点击确认收到之前，会不断提醒，因此需要更新下一次提醒时间为次日相同时间点
			r.changeNextTime(v.Id)
		}
	}
}
