package model

import "time"

type Post struct {
	Id        int
	CateId    int
	Type      int
	UserId    int
	Title     string
	Url       string
	Content   string
	Status    int
	ViewNum   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdatePost struct {
	Id        int
	CateId    *int
	Type      *int
	Title     *string
	Url       *string
	Content   *string
	Status    *int
	UpdatedAt *time.Time
}

type PostArchive struct {
	Ym    string
	Total string
}
