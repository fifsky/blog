package remind

import (
	"net/url"
	"strconv"
	"time"

	"app/provider/model"
	"github.com/goapt/golib/convert"
	"github.com/goapt/golib/robot"
	"github.com/goapt/logger"
	"github.com/ilibs/gosql/v2"

	"app/config"
	"app/pkg/aesutil"
)

func StartCron(conf *config.Config) {
	t := time.NewTicker(60 * time.Second)

	for t1 := range t.C {
		dingRemind(t1, conf)
	}
}

func numFormat(n int) string {
	if n < 10 {
		return "0" + strconv.Itoa(n)
	}
	return strconv.Itoa(n)
}

func message(content string, v *model.Reminds, conf *config.Config) {
	token, _ := aesutil.AesEncode(conf.Common.TokenSecret, convert.ToStr(v.Id))

	err := robot.CardMessage("⏰重要提醒⏰", content, []map[string]string{
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
		logger.Error(err)
	}
}

func changeNextTime(id int) {
	location, _ := time.LoadLocation("Asia/Shanghai")
	nextTime, _ := time.ParseInLocation("2006-01-02 15:04:00", time.Now().AddDate(0, 0, 1).Format("2006-01-02 15:04:00"), location)

	_, err := gosql.Model(&model.Reminds{Status: 2, NextTime: nextTime}).Where("id = ?", id).Update()

	if err != nil {
		logger.Error(err)
	}
}

func dingRemind(t time.Time, conf *config.Config) {
	reminds := make([]*model.Reminds, 0)
	err := gosql.Model(&reminds).All()
	if err != nil {
		logger.Error(err)
	}

	for _, v := range reminds {
		content := "提醒时间:" + time.Now().Format("2006-01-02 15:04:00") + " \n\n提醒内容:" + v.Content

		// 如果是等待确认的消息，则每天都需要提醒
		if v.Status == 2 {
			// 未确认的消息每天都需要在相同的时间点提醒
			if t.Format("15:04") == v.NextTime.Format("15:04") {
				message("🙋🏻‍再次提醒 \n "+content, v, conf)
				changeNextTime(v.Id)
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
			message(content, v, conf)
			changeNextTime(v.Id)
		}
	}
}
