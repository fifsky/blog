package model

import (
	"time"
)

type Users struct {
	Id        int       `form:"id" json:"id" db:"id"`
	Name      string    `form:"name" json:"name" db:"name"`
	Password  string    `form:"password" json:"password" db:"password"`
	NickName  string    `form:"nick_name" json:"nick_name" db:"nick_name"`
	Email     string    `form:"email" json:"email" db:"email"`
	Status    int       `form:"status" json:"status" db:"status"`
	Type      int       `form:"type" json:"type" db:"type"`
	CreatedAt time.Time `form:"-" json:"created_at" time_format:"2006-01-02 15:04:05" db:"created_at"`
	UpdatedAt time.Time `form:"-" json:"updated_at" time_format:"2006-01-02 15:04:05" db:"updated_at"`
}

func (u *Users) TableName() string {
	return "users"
}

func (u *Users) PK() string {
	return "id"
}
