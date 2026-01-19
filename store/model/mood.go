package model

import "time"

type Mood struct {
	Id        int
	Content   string
	UserId    int
	CreatedAt time.Time
}

type UserMood struct {
	Mood
	User *User
}

type UpdateMood struct {
	Id      int
	Content *string
}
