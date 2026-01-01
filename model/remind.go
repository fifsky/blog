package model

import "time"

type Remind struct {
	Id        int
	Type      int
	Content   string
	Month     int
	Week      int
	Day       int
	Hour      int
	Minute    int
	Status    int
	NextTime  time.Time
	CreatedAt time.Time
}

type UpdateRemind struct {
	Id       int
	Type     *int
	Content  *string
	Month    *int
	Week     *int
	Day      *int
	Hour     *int
	Minute   *int
	Status   *int
	NextTime *time.Time
}
