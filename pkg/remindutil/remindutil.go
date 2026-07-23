package remindutil

import (
	"time"

	"app/pkg/scheduler"
	"app/store/model"
)

// IsFixedDate 判断提醒规则是否为固定日期时间格式（如 "2006-01-02 15:04:00"），而非 cron 表达式。
// 固定日期任务为一次性提醒，确认后标记完成；cron 表达式为周期性提醒。
func IsFixedDate(cron string) bool {
	return len(cron) >= 10 && cron[4] == '-'
}

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

	// 固定日期格式（如 "2006-01-02 15:04:00"），直接解析为时间点
	if IsFixedDate(m.Cron) {
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
