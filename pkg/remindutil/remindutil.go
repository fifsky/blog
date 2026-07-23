package remindutil

import (
	"time"

	"app/pkg/scheduler"
	"app/store/model"
)

// NextTimeFromRule 根据提醒规则计算下一次触发时间
func NextTimeFromRule(from time.Time, m *model.Remind) time.Time {
	if m.Cron == "" {
		return time.Time{}
	}

	location, _ := time.LoadLocation("Asia/Shanghai")
	if location == nil {
		location = time.Local
	}
	base := from.In(location)

	// 如果是固定日期格式 "2006-01-02 15:04:00"
	if len(m.Cron) >= 10 && m.Cron[4] == '-' {
		t, err := time.ParseInLocation(time.DateTime, m.Cron, location)
		if err == nil {
			return t
		}
	}

	schedule, err := scheduler.ParseCronExpression(m.Cron)
	if err != nil {
		return time.Time{}
	}
	return schedule.Next(base)
}
