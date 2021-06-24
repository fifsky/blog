package model

import (
	"time"
)

type Comments struct {
	Id        int       `form:"id" json:"id" db:"id"`
	PostId    int       `form:"post_id" json:"post_id" db:"post_id"`
	Pid       int       `form:"pid" json:"pid" db:"pid"`
	Name      string    `form:"name" json:"name" db:"name"`
	Content   string    `form:"content" json:"content" db:"content"`
	IP        string    `form:"-" json:"ip" db:"ip"`
	CreatedAt time.Time `form:"-" json:"created_at" time_format:"2006-01-02 15:04:05" db:"created_at"`
}

func (c *Comments) TableName() string {
	return "comments"
}

func (c *Comments) PK() string {
	return "id"
}
