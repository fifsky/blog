package admin

import (
	"context"
	"testing"
	"time"

	"app/config"
	"app/model"
	adminv1 "app/proto/gen/admin/v1"
	"app/store"
	"app/testutil"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/goapt/dbunit"
)

func TestAdminUser_LoginUserGetListStatusCreate(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("users")...)
		svc := NewUser(store.New(db), &config.Config{})

		// LoginUser
		user := &model.User{
			Id:        1,
			Name:      "test",
			Password:  "test",
			NickName:  "test",
			Email:     "test@test.com",
			Status:    1,
			Type:      1,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		}
		ctx := SetLoginUser(context.Background(), user)
		resp, err := svc.LoginUser(ctx, &emptypb.Empty{})
		if err != nil || resp.Id == 0 {
			t.Fatalf("unexpected err=%v resp=%v", err, resp)
		}

		// Get
		resp2, err2 := svc.Get(context.Background(), &adminv1.GetUserRequest{Id: 1})
		if err2 != nil || resp2.Id == 0 {
			t.Fatalf("unexpected err=%v resp=%v", err2, resp2)
		}

		// List
		resp3, err3 := svc.List(context.Background(), &adminv1.UserListRequest{Page: 1})
		if err3 != nil || len(resp3.List) == 0 {
			t.Fatalf("unexpected err=%v list=%v", err3, resp3.List)
		}

		// Status
		_, err4 := svc.Status(context.Background(), &adminv1.UserStatusRequest{Id: 1})
		if err4 != nil {
			t.Fatalf("unexpected err=%v", err4)
		}
		_, err5 := svc.Status(context.Background(), &adminv1.UserStatusRequest{})
		if err5 == nil {
			t.Fatalf("expected validation error")
		}

		// Create
		_, err6 := svc.Create(context.Background(), &adminv1.UserCreateRequest{Name: "demo", Password: "123", NickName: "demo", Email: "demo@123.com", Type: 1})
		if err6 != nil {
			t.Fatalf("unexpected err=%v", err6)
		}
	})
}
