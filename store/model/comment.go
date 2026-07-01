package model

import "time"

// Comment 评论模型
type Comment struct {
	Id        int    // PK
	PostId    int    // 文章PID
	Pid       int    // 回复评论ID，顶层主评论为0
	Name      string // 名称
	Email     string // 邮箱（不对外展示，仅用于生成头像）
	Website   string // 网址
	Content   string // 内容
	ReplyName string // 被回复人昵称，回复的回复时用于平铺展示"A 回复 B"
	IP        string // IP
	CreatedAt time.Time
}

// CommentWithPost 评论关联文章信息，用于侧边栏最新评论和后台列表
type CommentWithPost struct {
	Comment
	PostTitle string // 关联文章标题
	PostUrl   string // 关联文章缩略名
}
