package model

import "time"

type Photo struct {
	Id          int
	Title       string
	Description string
	Src         string
	Thumbnail   string
	Province    string
	City        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UpdatePhoto struct {
	Id          int
	Title       *string
	Description *string
	Province    *string
	City        *string
}
