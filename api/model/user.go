package model

import "time"

type User struct {
	Id        int
	Name      string
	Password  string
	NickName  string
	Email     string
	Status    int
	Type      int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateUser struct {
	Id        int
	Name      *string
	Password  *string
	NickName  *string
	Email     *string
	Status    *int
	Type      *int
	UpdatedAt *time.Time
}
