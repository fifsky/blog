package model

import (
	"time"
)

type Posts struct {
	Id        int       `json:"id" db:"id"`
	CateId    int       `json:"cate_id" db:"cate_id" binding:"required"`
	Type      int       `json:"type" db:"type" binding:"required"`
	UserId    int       `json:"user_id" db:"user_id"`
	Title     string    `json:"title" db:"title" binding:"required"`
	Url       string    `json:"url" db:"url"`
	Content   string    `json:"content" db:"content" binding:"required"`
	Status    int       `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" time_format:"2006-01-02 15:04:05" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" time_format:"2006-01-02 15:04:05" db:"updated_at"`
}

func (p *Posts) TableName() string {
	return "posts"
}

func (p *Posts) PK() string {
	return "id"
}
