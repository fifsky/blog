package model

import (
	"time"
)

type Guestbook struct {
	Id        int
	Name      string
	Content   string
	Ip        string
	Top       int
	CreatedAt time.Time
}
