package model

import "time"

// 链接状态常量
const (
	LinkStatusPending  = "PENDING"  // 审核中
	LinkStatusApproved = "APPROVED" // 审核通过
)

type Link struct {
	Id        int
	Name      string
	Url       string
	Desc      string
	Status    string
	CreatedAt time.Time
}

type UpdateLink struct {
	Id     int
	Name   *string
	Url    *string
	Desc   *string
	Status *string
}
