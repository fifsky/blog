package model

import "time"

// Mood 心情模型
type Mood struct {
	Id        int    // PK
	Content   string // 心情内容
	UserId    int    // 用户ID
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserMood 心情及关联用户
type UserMood struct {
	Mood
	User *User
}

// UpdateMood 更新心情参数
type UpdateMood struct {
	Id      int
	Content *string
}
