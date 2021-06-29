package model

import (
	"time"
)

type Comments struct {
	Id        int       `json:"id" db:"id"`
	PostId    int       `json:"post_id" db:"post_id" binding:"required"`
	Pid       int       `json:"pid" db:"pid"`
	Name      string    `json:"name" db:"name" binding:"required"`
	Content   string    `json:"content" db:"content" binding:"required"`
	IP        string    `json:"ip" db:"ip"`
	CreatedAt time.Time `json:"created_at" time_format:"2006-01-02 15:04:05" db:"created_at"`
}

func (c *Comments) TableName() string {
	return "comments"
}

func (c *Comments) PK() string {
	return "id"
}
