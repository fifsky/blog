package model

import "time"

type Comment struct {
	Id        int
	PostId    int
	Pid       int
	Name      string
	Content   string
	IP        string
	CreatedAt time.Time
}
type NewComment struct {
	Comment
	Type         int
	ArticleTitle string
	Url          string
}
