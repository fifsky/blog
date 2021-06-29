package model

import (
	"time"
)

type Cates struct {
	Id        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" binding:"required"`
	Desc      string    `json:"desc" db:"desc"`
	Domain    string    `json:"domain" db:"domain" binding:"required"`
	CreatedAt time.Time `json:"created_at" time_format:"2006-01-02 15:04:05" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" time_format:"2006-01-02 15:04:05" db:"updated_at"`
}

func (c *Cates) TableName() string {
	return "cates"
}

func (c *Cates) PK() string {
	return "id"
}
