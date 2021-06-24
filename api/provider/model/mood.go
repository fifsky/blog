package model

import (
	"time"
)

type Moods struct {
	Id        int       `form:"id" json:"id" db:"id"`
	Content   string    `form:"content" json:"content" db:"content"`
	UserId    int       `form:"user_id" json:"user_id" db:"user_id"`
	CreatedAt time.Time `form:"-" json:"created_at" time_format:"2006-01-02 15:04:05" db:"created_at"`
}

func (m *Moods) TableName() string {
	return "moods"
}

func (m *Moods) PK() string {
	return "id"
}
