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

func numFormat(n int) string {
	if n < 10 {
		return "0" + strconv.Itoa(n)
	}
	return strconv.Itoa(n)
}

func NextTimeFromRule(from time.Time, m *model.Remind) time.Time {
	location, _ := time.LoadLocation("Asia/Shanghai")
	if location == nil {
		location = time.Local
	}
	base := from.In(location)

	switch m.Type {
	case 0:
		year := base.Year()
		if !m.CreatedAt.IsZero() {
			year = m.CreatedAt.In(location).Year()
		}
		month := time.Month(m.Month)
		if month == 0 {
			month = base.Month()
		}
		day := m.Day
		if day == 0 {
			day = base.Day()
		}
		return time.Date(year, month, day, m.Hour, m.Minute, 0, 0, location)
	case 1:
		return base.Add(time.Minute)
	case 2:
		return base.Add(time.Hour)
	case 3:
		week := m.Week
		if week < 1 || week > 7 {
			week = int(base.Weekday()) + 1
		}
		target := time.Weekday((week - 1 + 7) % 7)
		candidate := time.Date(base.Year(), base.Month(), base.Day(), m.Hour, m.Minute, 0, 0, location)
		for candidate.Weekday() != target {
			candidate = candidate.AddDate(0, 0, 1)
		}
		if !candidate.After(base) {
			candidate = candidate.AddDate(0, 0, 7)
		}
		return candidate
	case 4:
		candidate := time.Date(base.Year(), base.Month(), base.Day(), m.Hour, m.Minute, 0, 0, location)
		if !candidate.After(base) {
			candidate = candidate.AddDate(0, 0, 1)
		}
		return candidate
	case 5:
		day := m.Day
		if day <= 0 {
			day = base.Day()
		}
		candidate := time.Date(base.Year(), base.Month(), day, m.Hour, m.Minute, 0, 0, location)
		if !candidate.After(base) {
			candidate = candidate.AddDate(0, 1, 0)
		}
		return candidate
	case 6:
		month := time.Month(m.Month)
		if month == 0 {
			month = base.Month()
		}
		day := m.Day
		if day == 0 {
			day = base.Day()
		}
		candidate := time.Date(base.Year(), month, day, m.Hour, m.Minute, 0, 0, location)
		if !candidate.After(base) {
			candidate = candidate.AddDate(1, 0, 0)
		}
		return candidate
	default:
		return base
	}
}

func (r *Remind) buildMessage(title, content string, v *model.Remind) messenger.Message {
	token, _ := aesutil.AesEncode(r.conf.Common.TokenSecret, strconv.Itoa(v.Id))

	changeURL := "https://api.fifsky.com/blog/remind/change?token=" + url.QueryEscape(token)
	delayURL := "https://api.fifsky.com/blog/remind/delay?token=" + url.QueryEscape(token)

	return messenger.Message{
		Title:   title,
		Content: content,
		Time:    time.Now().Format("2006-01-02 15:04"),
		Actions: []messenger.Action{
			{Title: "收到提醒", URL: changeURL},
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

		isRemind := false
		switch v.Type {
		case 0: // 固定时间
			if t.Format("2006-01-02 15:04:00") == v.NextTime.Format("2006-01-02 15:04:00") {
				isRemind = true
			}
		case 1: // 每分钟
			isRemind = true
		case 2: // 每小时
			if t.Format("04:00") == numFormat(v.Hour)+":00" {
				isRemind = true
			}
		case 3: // 每周
			if t.Weekday().String() == time.Weekday(v.Week-1).String() && t.Format("15:04:00") == numFormat(v.Hour)+":"+numFormat(v.Minute)+":00" {
				isRemind = true
			}
		case 4: // 每天
			if t.Format("15:04:00") == numFormat(v.Hour)+":"+numFormat(v.Minute)+":00" {
				isRemind = true
			}
		case 5: // 每月
			if t.Format("02 15:04:00") == numFormat(v.Day)+" "+numFormat(v.Hour)+":"+numFormat(v.Minute)+":00" {
				isRemind = true
			}
		case 6: // 每年
			if t.Format("01-02 15:04:00") == numFormat(v.Month)+"-"+numFormat(v.Day)+" "+numFormat(v.Hour)+":"+numFormat(v.Minute)+":00" {
				isRemind = true
			}
		}

		if isRemind {
			v2 := v
			r.message("⏰重要提醒⏰", content, &v2)
			// 如果发出提醒，在用户没有点击确认收到之前，会不断提醒，因此需要更新下一次提醒时间为次日相同时间点
			r.changeNextTime(v.Id)
		}
	}
}
