package model

import "time"

type Remind struct {
	Id        int
	Cron      string
	Content   string
	Status    int
	NextTime  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateRemind struct {
	Id       int
	Cron     *string
	Content  *string
	Status   *int
	NextTime *time.Time
}
