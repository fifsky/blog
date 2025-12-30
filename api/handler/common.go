package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"app/model"
	"app/pkg/validate"
)

func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}

	// 仅校验结构体类型，避免 map/slice 等类型触发无意义校验
	if reflect.TypeOf(v).Kind() == reflect.Struct {
		if err := validate.Validate(&v); err != nil {
			return v, err
		}
	}
	return v, nil
}

type loginUserKey struct{}

func SetLoginUser(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, loginUserKey{}, user)
}

func getLoginUser(ctx context.Context) *model.User {
	u := ctx.Value(loginUserKey{})
	if u == nil {
		panic(errors.New("login user not found"))
	}
	return u.(*model.User)
}

type remindKey struct{}

func SetRemind(ctx context.Context, remind *model.Remind) context.Context {
	return context.WithValue(ctx, remindKey{}, remind)
}

func getRemind(ctx context.Context) *model.Remind {
	r := ctx.Value(remindKey{})
	if r == nil {
		panic(errors.New("remind not found"))
	}
	return r.(*model.Remind)
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
