package model

import "time"

// AuthProfile 用户身份验证特征
type AuthProfile struct {
	Id                    int
	UserId                int
	IdentityDescription   string
	VerificationThreshold float64
	MaxAttempts           int
	CreatedAt             time.Time
	UpdatedAt             time.Time
}
