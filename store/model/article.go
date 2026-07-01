package model

import "time"

// Post 文章模型
type Post struct {
	Id        int    // PK
	CateId    int    // 文章分类ID
	Type      int    // 类型：1:文章,2:页面
	UserId    int    // 文章作者ID
	Title     string // 文章标题
	Url       string // 页面缩略名
	Content   string // 文章内容
	Tags      Tags   // 文章标签
	Status    int    // 状态，1正常 2删除 3草稿
	ViewNum   int    // 浏览次数
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
	Status    *int
	UpdatedAt *time.Time
}

// PostArchive 文章归档统计
type PostArchive struct {
	Ym    string // 年月，格式 YYYY/MM
	Total string // 文章数量
}
