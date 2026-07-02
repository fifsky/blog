package model

import "time"

// PostStatus 文章状态类型
type PostStatus string

// 文章状态常量
const (
	PostStatusActive  PostStatus = "ACTIVE"  // 正常
	PostStatusDeleted PostStatus = "DELETED" // 删除
	PostStatusDraft   PostStatus = "DRAFT"   // 草稿
)

// Post 文章模型
type Post struct {
	Id        int        // PK
	CateId    int        // 文章分类ID
	Type      int        // 类型：1:文章,2:页面
	UserId    int        // 文章作者ID
	Title     string     // 文章标题
	Url       string     // 页面缩略名
	Content   string     // 文章内容
	Tags      Tags       // 文章标签
	Status    PostStatus // 状态:ACTIVE正常,DELETED删除,DRAFT草稿
	ViewNum   int        // 浏览次数
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UpdatePost 更新文章参数
type UpdatePost struct {
	Id        int
	CateId    *int
	Type      *int
	Title     *string
	Url       *string
	Content   *string
	Tags      *Tags
	Status    *PostStatus
	UpdatedAt *time.Time
}

// PostArchive 文章归档统计
type PostArchive struct {
	Ym    string // 年月，格式 YYYY/MM
	Total string // 文章数量
}
