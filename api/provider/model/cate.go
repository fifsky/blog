package model

import (
	"time"
)

type Cates struct {
	Id        int       `form:"id" json:"id" db:"id"`
	Name      string    `form:"name" json:"name" db:"name"`
	Desc      string    `form:"desc" json:"desc" db:"desc"`
	Domain    string    `form:"domain" json:"domain" db:"domain"`
	CreatedAt time.Time `form:"-" json:"created_at" time_format:"2006-01-02 15:04:05" db:"created_at"`
	UpdatedAt time.Time `form:"-" json:"updated_at" time_format:"2006-01-02 15:04:05" db:"updated_at"`
}

func (c *Cates) TableName() string {
	return "cates"
}

func (c *Cates) PK() string {
	return "id"
}
