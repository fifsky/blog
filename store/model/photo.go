package model

import "time"

type Photo struct {
	Id          int
	Title       string
	Description string
	Src         string
	Thumbnail   string
	Province    int
	City        int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UpdatePhoto struct {
	Id          int
	Title       *string
	Description *string
	Province    *int
	City        *int
}
