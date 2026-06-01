package model

import "time"

type User struct {
	Id         int
	Name       string
	Password   string
	NickName   string
	Email      string
	Status     int
	Type       int
	TotpSecret string
	Openid     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type UpdateUser struct {
	Id         int
	Name       *string
	Password   *string
	NickName   *string
	Email      *string
	Status     *int
	Type       *int
	TotpSecret *string
	UpdatedAt  *time.Time
}
