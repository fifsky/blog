package model

import (
	"time"

	"github.com/ilibs/gosql/v2"
)

type Reminds struct {
	Id        int       `form:"id" json:"id" db:"id"`
	Type      int       `form:"type" json:"type" db:"type"`
	Content   string    `form:"content" json:"content" db:"content"`
	Month     int       `form:"month" json:"month" db:"month"`
	Week      int       `form:"week" json:"week" db:"week"`
	Day       int       `form:"day" json:"day" db:"day"`
	Hour      int       `form:"hour" json:"hour" db:"hour"`
	Minute    int       `form:"minute" json:"minute" db:"minute"`
	Status    int       `form:"-" json:"status" db:"status"`
	NextTime  time.Time `form:"-" json:"next_time" db:"next_time"`
	CreatedAt time.Time `form:"-" json:"created_at" time_format:"2006-01-02 15:04:05" db:"created_at"`
}

func (r *Reminds) TableName() string {
	return "reminds"
}

func (r *Reminds) PK() string {
	return "id"
}

func RemindGetList(start int, num int) ([]*Reminds, error) {
	var m = make([]*Reminds, 0)
	start = (start - 1) * num
	err := gosql.Model(&m).OrderBy("id desc").Limit(num).Offset(start).All()
	if err != nil {
		return nil, err
	}
	return m, nil
}
