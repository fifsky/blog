package model

import (
	"time"
)

// Guestbook 留言板模型
type Guestbook struct {
	Id        int    // PK
	Name      string // 昵称
	Content   string // 内容
	Ip        string // IP
	Top       int    // 1置顶
	CreatedAt time.Time
}
