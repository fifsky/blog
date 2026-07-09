package admin

import (
	"context"
	"testing"
	"time"

	"app/config"
	"app/pkg/dbunit"
	adminv1 "app/proto/gen/admin/v1"
	"app/store"
	"app/store/model"
	"app/testutil"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/stretchr/testify/assert"
)

func TestAdminUser_LoginUserGetListStatusCreate(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("users"))
		svc := NewUser(store.New(db), &config.Config{})

		t.Run("获取当前登录用户", func(t *testing.T) {
			user := &model.User{
				Id:        1,
				Name:      "test",
				Password:  "test",
				NickName:  "test",
				Email:     "test@test.com",
				Status:    model.UserStatusActive,
				Type:      1,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			}
			ctx := SetLoginUser(context.Background(), user)
			resp, err := svc.LoginUser(ctx, &emptypb.Empty{})
			assert.NoError(t, err)
			assert.NotZero(t, resp.GetId())
		})

		t.Run("获取用户详情", func(t *testing.T) {
			resp, err := svc.Get(context.Background(), adminv1.GetUserRequest_builder{Id: 1}.Build())
			assert.NoError(t, err)
			assert.NotZero(t, resp.GetId())
		})

		t.Run("获取用户列表", func(t *testing.T) {
			resp, err := svc.List(context.Background(), adminv1.UserListRequest_builder{Page: 1}.Build())
			assert.NoError(t, err)
			assert.NotEmpty(t, resp.GetList())
		})

		t.Run("更新用户状态", func(t *testing.T) {
			// 正常更新
			_, err := svc.Status(context.Background(), adminv1.UserStatusRequest_builder{Id: 1}.Build())
			assert.NoError(t, err)

			// 参数验证失败
			_, err = svc.Status(context.Background(), adminv1.UserStatusRequest_builder{}.Build())
			assert.Error(t, err)
		})

		t.Run("创建用户", func(t *testing.T) {
			_, err := svc.Create(context.Background(), adminv1.UserCreateRequest_builder{Name: "demo",
				Password: "123",
				NickName: "demo",
				Email:    "demo@123.com",
				Type:     1}.Build(),
			)
			assert.NoError(t, err)
		})
	})
}
