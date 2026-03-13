package model

import "time"

// AuthSessionStatus 会话状态类型
type AuthSessionStatus string

const (
	AuthSessionActive  AuthSessionStatus = "active"
	AuthSessionSuccess AuthSessionStatus = "success"
	AuthSessionFailed  AuthSessionStatus = "failed"
	AuthSessionExpired AuthSessionStatus = "expired"
)

// AuthSession AI验证会话
type AuthSession struct {
	Id            int
	SessionId     string
	UserId        int
	AttemptCount  int
	VerifiedScore float64
	Status        AuthSessionStatus
	ExpiresAt     time.Time
	CreatedAt     time.Time
}
