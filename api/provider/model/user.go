package model

import (
	"time"
)

type Users struct {
	Id        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" binding:"required"`
	Password  string    `json:"password" db:"password"`
	NickName  string    `json:"nick_name" db:"nick_name" binding:"required"`
	Email     string    `json:"email" db:"email"`
	Status    int       `json:"status" db:"status"`
	Type      int       `json:"type" db:"type" binding:"required"`
	CreatedAt time.Time `json:"created_at" time_format:"2006-01-02 15:04:05" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" time_format:"2006-01-02 15:04:05" db:"updated_at"`
}

func (u *Users) TableName() string {
	return "users"
}

func (u *Users) PK() string {
	return "id"
}
