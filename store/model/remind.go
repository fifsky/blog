package model

import "time"

// RemindStatus 提醒状态类型
type RemindStatus string

// 提醒状态常量
const (
	RemindStatusActive  RemindStatus = "ACTIVE"  // 正常
	RemindStatusPending RemindStatus = "PENDING" // 等待确认
	RemindStatusDone    RemindStatus = "DONE"    // 已完成
)

// Remind 提醒模型
type Remind struct {
	Id        int          // PK
	Cron      string       // cron表达式或固定时间
	Content   string       // 提醒内容
	Status    RemindStatus // 状态:ACTIVE正常,PENDING等待确认,DONE已完成
	NextTime  time.Time    // 下次提醒时间
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UpdateRemind 更新提醒参数
type UpdateRemind struct {
	Id       int
	Cron     *string
	Content  *string
	Status   *RemindStatus
	NextTime *time.Time
}
