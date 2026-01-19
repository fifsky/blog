package model

import "time"

type Link struct {
	Id        int
	Name      string
	Url       string
	Desc      string
	CreatedAt time.Time
}

type UpdateLink struct {
	Id   int
	Name *string
	Url  *string
	Desc *string
}
