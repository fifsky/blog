package model

import "time"

type Cate struct {
	Id        int
	Name      string
	Desc      string
	Domain    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CateArtivleCount struct {
	Id        int
	Name      string
	Desc      string
	Domain    string
	CreatedAt time.Time
	UpdatedAt time.Time
	Num       int
}

type UpdateCate struct {
	Id        int
	Name      *string
	Desc      *string
	Domain    *string
	UpdatedAt *time.Time
}
