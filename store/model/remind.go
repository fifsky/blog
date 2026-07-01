package model

import "time"

// Remind 提醒模型
type Remind struct {
	Id        int       // PK
	Cron      string    // cron表达式或固定时间
	Content   string    // 提醒内容
	Status    int       // 状态
	NextTime  time.Time // 下次提醒时间
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UpdateRemind 更新提醒参数
type UpdateRemind struct {
	Id       int
	Cron     *string
	Content  *string
	Status   *int
	NextTime *time.Time
}
