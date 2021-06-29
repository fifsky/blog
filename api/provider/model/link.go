package model

import (
	"time"
)

type Links struct {
	Id        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" binding:"required"`
	Url       string    `json:"url" db:"url" binding:"required"`
	Desc      string    `json:"desc" db:"desc"`
	CreatedAt time.Time `json:"created_at" time_format:"2006-01-02 15:04:05" db:"created_at"`
}

func (l *Links) TableName() string {
	return "links"
}

func (l *Links) PK() string {
	return "id"
}
