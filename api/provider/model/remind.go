package model

import (
	"time"
)

type Reminds struct {
	Id        int       `json:"id" db:"id"`
	Type      int       `json:"type" db:"type" binding:"required"` // 0固定，1每分钟，2每个小时，3每周，4，每天，5，每月，6，每年
	Content   string    `json:"content" db:"content" binding:"required"`
	Month     int       `json:"month" db:"month"`
	Week      int       `json:"week" db:"week"`
	Day       int       `json:"day" db:"day"`
	Hour      int       `json:"hour" db:"hour"`
	Minute    int       `json:"minute" db:"minute"`
	Status    int       `json:"status" db:"status"`
	NextTime  time.Time `json:"next_time" db:"next_time"`
	CreatedAt time.Time `json:"created_at" time_format:"2006-01-02 15:04:05" db:"created_at"`
}

func (r *Reminds) TableName() string {
	return "reminds"
}

func (r *Reminds) PK() string {
	return "id"
}
