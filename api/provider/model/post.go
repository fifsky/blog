package model

import (
	"time"
)

type Posts struct {
	Id        int       `form:"id" json:"id" db:"id"`
	CateId    int       `form:"cate_id" json:"cate_id" db:"cate_id"`
	Type      int       `form:"type" json:"type" db:"type"`
	UserId    int       `form:"user_id" json:"user_id" db:"user_id"`
	Title     string    `form:"title" json:"title" db:"title"`
	Url       string    `form:"url" json:"url" db:"url"`
	Content   string    `form:"content" json:"content" db:"content"`
	Status    int       `form:"-" json:"status" db:"status"`
	CreatedAt time.Time `form:"-" json:"created_at" time_format:"2006-01-02 15:04:05" db:"created_at"`
	UpdatedAt time.Time `form:"-" json:"updated_at" time_format:"2006-01-02 15:04:05" db:"updated_at"`
}

func (p *Posts) TableName() string {
	return "posts"
}

func (p *Posts) PK() string {
	return "id"
}
