package model

import (
	"time"
)

type Links struct {
	Id        int       `form:"id" json:"id" db:"id"`
	Name      string    `form:"name" json:"name" db:"name"`
	Url       string    `form:"url" json:"url" db:"url"`
	Desc      string    `form:"desc" json:"desc" db:"desc"`
	CreatedAt time.Time `form:"-" json:"created_at" time_format:"2006-01-02 15:04:05" db:"created_at"`
}

func (l *Links) TableName() string {
	return "links"
}

func (l *Links) PK() string {
	return "id"
}
