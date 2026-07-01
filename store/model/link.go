package model

import "time"

// 链接状态常量
const (
	LinkStatusPending  = "PENDING"  // 审核中
	LinkStatusApproved = "APPROVED" // 审核通过
)

// Link 友情链接模型
type Link struct {
	Id        int    // PK
	Name      string // 链接名称
	Url       string // 链接地址
	Desc      string // 链接描述
	Status    string // 状态:PENDING审核中,APPROVED审核通过
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UpdateLink 更新链接参数
type UpdateLink struct {
	Id     int
	Name   *string
	Url    *string
	Desc   *string
	Status *string
}
