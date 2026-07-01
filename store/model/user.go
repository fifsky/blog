package model

import "time"

// User 用户模型
type User struct {
	Id         int    // PK
	Name       string // 用户名
	Password   string // 密码
	NickName   string // 昵称
	Email      string // 邮箱
	Status     int    // 状态：1正常，2删除
	Type       int    // 类型：1:管理员,2:编辑
	TotpSecret string // TOTP密钥
	Openid     string // 小程序openid
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// UpdateUser 更新用户参数
type UpdateUser struct {
	Id         int
	Name       *string
	Password   *string
	NickName   *string
	Email      *string
	Status     *int
	Type       *int
	TotpSecret *string
	UpdatedAt  *time.Time
}
