package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"app/model"
	"app/pkg/validate"
)

func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}

	if err := validate.Validate(&v); err != nil {
		return v, err
	}
	return v, nil
}

func getLoginUser(ctx context.Context) *model.User {
	// 从请求上下文取出登录用户信息
	if u := ctx.Value("userInfo"); u != nil {
		return u.(*model.User)
	}
	return nil
}

func getRemind(ctx context.Context) *model.Remind {
	// 从请求上下文取出登录用户信息
	if u := ctx.Value("remind"); u != nil {
		return u.(*model.Remind)
	}
	return nil
}

// 获取客户端IP
func clientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		parts := strings.Split(fwd, ",")
		return strings.TrimSpace(parts[0])
	}
	return r.RemoteAddr
}

func totalPages(total, pagingNum int) int {
	if total == 0 {
		return 1
	}
	if total%pagingNum == 0 {
		return total / pagingNum
	}
	return total/pagingNum + 1
}

func parseTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	// 尝试标准格式
	if t, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
		return t, nil
	}
	// 尝试RFC3339
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("invalid time format: %s", s)
}

// 统一在 request_types 中使用结构体接收参数，移除此类 map 辅助方法
