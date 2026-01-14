package remind

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"app/config"
	"app/model"
	"app/pkg/aesutil"
	"github.com/goapt/logger"
	"app/pkg/wechat"
	"app/store"
)

type Remind struct {
	store *store.Store
	conf  *config.Config
	robot *wechat.Robot
}

func New(s *store.Store, conf *config.Config, robot *wechat.Robot) *Remind {
	return &Remind{
		store: s,
		conf:  conf,
		robot: robot,
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

type barkRequest struct {
	Body     string `json:"body,omitempty"`
	Title    string `json:"title"`
	Badge    int    `json:"badge"`
	Url      string `json:"url,omitempty"`
	Markdown string `json:"markdown"`
}

func (r *Remind) messageForBark(content string, v *model.Remind) {
	token, _ := aesutil.AesEncode(r.conf.Common.TokenSecret, strconv.Itoa(v.Id))

	markdown := `
%s

[收到提醒](%s)  [稍后提醒](%s)
`

	body := barkRequest{
		Title:    "⏰重要提醒⏰",
		Badge:    1,
		Markdown: fmt.Sprintf(markdown, content, "https://api.fifsky.com/api/remind/change?token="+url.QueryEscape(token), "https://api.fifsky.com/api/remind/delay?token="+url.QueryEscape(token)),
	}

	reqBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", r.conf.Common.NotifyUrl, bytes.NewReader(reqBody))

	if err != nil {
		logger.Default().Error("remind request bark error", slog.String("err", err.Error()))
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Basic "+r.conf.Common.NotifyToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Default().Error("remind request bark error", slog.String("err", err.Error()))
	}
	defer resp.Body.Close()
}

func (r *Remind) message(content string, v *model.Remind) {
	token, _ := aesutil.AesEncode(r.conf.Common.TokenSecret, strconv.Itoa(v.Id))

	err := r.robot.CardMessage("⏰重要提醒⏰", content, []map[string]string{
		{
			"title":     "收到提醒",
			"actionURL": "https://api.fifsky.com/api/remind/change?token=" + url.QueryEscape(token),
		},
		{
			"title":     "稍后提醒",
			"actionURL": "https://api.fifsky.com/api/remind/delay?token=" + url.QueryEscape(token),
		},
	})
	if err != nil {
		logger.Default().Error("remind request robot error", slog.String("err", err.Error()))
	}
	r.messageForBark(content, v)
}

func (r *Remind) changeNextTime(id int) {
	location, _ := time.LoadLocation("Asia/Shanghai")
	nextTime, _ := time.ParseInLocation("2006-01-02 15:04:00", time.Now().AddDate(0, 0, 1).Format("2006-01-02 15:04:00"), location)

	_ = r.store.UpdateRemindStatus(context.Background(), id, 2)
	_ = r.store.UpdateRemindNextTime(context.Background(), id, nextTime)
}

func (r *Remind) run(t time.Time) {
	reminds, _ := r.store.RemindAll(context.Background())

	for _, v := range reminds {
		content := "提醒时间:" + time.Now().Format("2006-01-02 15:04:00") + " \n\n提醒内容:" + v.Content

		// 如果是等待确认的消息，则每天都需要提醒
		if v.Status == 2 {
			// 未确认的消息每天都需要在相同的时间点提醒
			if t.Format("15:04") == v.NextTime.Format("15:04") {
				v2 := v
				r.message("🙋🏻‍再次提醒 \n "+content, &v2)
				r.changeNextTime(v.Id)
			}

			continue
		}

		isRemind := false
		switch v.Type {
		case 0: // 固定时间
			if t.Format("2006-01-02 15:04:00") == v.CreatedAt.Format("2006")+"-"+numFormat(v.Month)+"-"+numFormat(v.Day)+" "+numFormat(v.Hour)+":"+numFormat(v.Minute)+":00" {
				isRemind = true
			}
		case 1: // 每分钟
			isRemind = true
		case 2: // 每小时
			if t.Format("04:00") == numFormat(v.Hour)+":00" {
				isRemind = true
			}
		case 3: // 每天
			if t.Format("15:04:00") == numFormat(v.Hour)+":"+numFormat(v.Minute)+":00" {
				isRemind = true
			}
		case 4: // 每周
			if t.Weekday().String() == time.Weekday(v.Week-1).String() && t.Format("15:04:00") == numFormat(v.Hour)+":"+numFormat(v.Minute)+":00" {
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
			r.message(content, &v2)
			r.changeNextTime(v.Id)
		}
	}
}
