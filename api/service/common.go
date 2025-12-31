package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"app/model"
	"app/response"
	"buf.build/go/protovalidate"
	"google.golang.org/protobuf/proto"
)

type ServiceFunc[T any, R any] func(ctx context.Context, r T) (R, error)

func Wrap[T any, R any](f ServiceFunc[T, R]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx  = r.Context()
			data T
			err  error
		)

		if data, err = decode[T](r); err != nil {
			response.Fail(w, 201, fmt.Sprintf("参数错误：%s", err))
			return
		}

		var result R
		if result, err = f(ctx, data); err != nil {
			response.Fail(w, 202, err.Error())
			return
		}

		response.Success(w, result)
	}
}

func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}

	if msg, ok := any(v).(proto.Message); ok {
		if err := protovalidate.Validate(msg); err != nil {
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
