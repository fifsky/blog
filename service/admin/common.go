package admin

import (
	"context"
	"errors"

	"app/store/model"
)

type loginUserKey struct{}

func SetLoginUser(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, loginUserKey{}, user)
}

func GetLoginUser(ctx context.Context) *model.User {
	u := ctx.Value(loginUserKey{})
	if u == nil {
		panic(errors.New("login user not found"))
	}
	return u.(*model.User)
}
