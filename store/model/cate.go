package model

import "time"

// Cate 文章分类模型
type Cate struct {
	Id        int // PK
	Name      string
	Desc      string // 分类描述
	Domain    string // 分类域名path
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CateArtivleCount 分类及文章数量统计
type CateArtivleCount struct {
	Id        int
	Name      string
	Desc      string // 分类描述
	Domain    string // 分类域名path
	CreatedAt time.Time
	UpdatedAt time.Time
	Num       int // 分类下文章数量
}

// UpdateCate 更新分类参数
type UpdateCate struct {
	Id        int
	Name      *string
	Desc      *string
	Domain    *string
	UpdatedAt *time.Time
}
