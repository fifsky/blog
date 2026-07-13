package store

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"app/pkg/dbunit"
	"app/store/model"
	"app/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_GetList(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("users"))
		s := New(db)
		users, err := s.ListUser(context.Background(), 1, 1)
		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.True(t, len(users) > 0)
	})
}

func TestUser_GetUser(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("users"))
		s := New(db)
		ret, err := s.GetUser(context.Background(), 1)
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		assert.Equal(t, "test", ret.Name)
	})
}

func TestUser_GetUser_NotFound(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("users"))
		s := New(db)
		ret, err := s.GetUser(context.Background(), 999)
		assert.Error(t, err)
		assert.Nil(t, ret)
		assert.True(t, errors.Is(err, sql.ErrNoRows))
	})
}

func TestUser_CountUserTotal(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("users"))
		s := New(db)
		total, err := s.CountUserTotal(context.Background())
		require.NoError(t, err)
		assert.Equal(t, 3, total)
	})

	t.Run("空表", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), nil)
			s := New(db)
			total, err := s.CountUserTotal(context.Background())
			require.NoError(t, err)
			assert.Equal(t, 0, total)
		})
	})
}

func TestUser_GetUserByName(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
		wantId   int
	}{
		{name: "存在用户test", username: "test", wantErr: false, wantId: 1},
		{name: "存在用户rita", username: "rita", wantErr: false, wantId: 2},
		{name: "不存在用户", username: "nobody", wantErr: true, wantId: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("users"))
				s := New(db)
				ret, err := s.GetUserByName(context.Background(), tt.username)
				if tt.wantErr {
					require.Error(t, err)
					assert.Nil(t, ret)
					return
				}
				require.NoError(t, err)
				assert.NotNil(t, ret)
				assert.Equal(t, tt.wantId, ret.Id)
				assert.Equal(t, tt.username, ret.Name)
			})
		})
	}
}

func TestUser_CreateUser(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("users"))
		s := New(db)

		beforeTotal, err := s.CountUserTotal(context.Background())
		require.NoError(t, err)

		now := time.Now()
		u := &model.User{
			Name:      "newuser",
			Password:  "password123",
			NickName:  "新用户",
			Email:     "new@example.com",
			Status:    model.UserStatusActive,
			Type:      2,
			CreatedAt: now,
			UpdatedAt: now,
		}

		id, err := s.CreateUser(context.Background(), u)
		require.NoError(t, err)
		assert.Greater(t, id, int64(0))

		afterTotal, err := s.CountUserTotal(context.Background())
		require.NoError(t, err)
		assert.Equal(t, beforeTotal+1, afterTotal)

		// 验证创建的用户数据
		got, err := s.GetUserByName(context.Background(), "newuser")
		require.NoError(t, err)
		assert.Equal(t, "新用户", got.NickName)
		assert.Equal(t, "new@example.com", got.Email)
		assert.Equal(t, model.UserStatusActive, got.Status)
		assert.Equal(t, 2, got.Type)
	})

	t.Run("重复用户名报错", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("users"))
			s := New(db)

			now := time.Now()
			u := &model.User{
				Name:      "test", // 已存在
				Password:  "password",
				NickName:  "test",
				Email:     "test2@example.com",
				Status:    model.UserStatusActive,
				Type:      1,
				CreatedAt: now,
				UpdatedAt: now,
			}
			_, err := s.CreateUser(context.Background(), u)
			require.Error(t, err)
		})
	})
}

func TestUser_UpdateUser(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		update  *model.UpdateUser
		wantErr bool
		check   func(t *testing.T, s *Store, id int)
	}{
		{
			name: "更新昵称",
			id:   1,
			update: &model.UpdateUser{
				Id:       1,
				NickName: new("新昵称"),
			},
			wantErr: false,
			check: func(t *testing.T, s *Store, id int) {
				u, err := s.GetUser(context.Background(), id)
				require.NoError(t, err)
				assert.Equal(t, "新昵称", u.NickName)
			},
		},
		{
			name: "更新密码",
			id:   1,
			update: &model.UpdateUser{
				Id:       1,
				Password: new("newpassword"),
			},
			wantErr: false,
			check: func(t *testing.T, s *Store, id int) {
				u, err := s.GetUser(context.Background(), id)
				require.NoError(t, err)
				assert.Equal(t, "newpassword", u.Password)
			},
		},
		{
			name: "更新状态为DELETED",
			id:   1,
			update: &model.UpdateUser{
				Id:     1,
				Status: new(model.UserStatusDeleted),
			},
			wantErr: false,
			check: func(t *testing.T, s *Store, id int) {
				u, err := s.GetUser(context.Background(), id)
				require.NoError(t, err)
				assert.Equal(t, model.UserStatusDeleted, u.Status)
			},
		},
		{
			name: "更新多个字段",
			id:   1,
			update: &model.UpdateUser{
				Id:       1,
				NickName: new("multi"),
				Email:    new("multi@example.com"),
				Type:     new(2),
			},
			wantErr: false,
			check: func(t *testing.T, s *Store, id int) {
				u, err := s.GetUser(context.Background(), id)
				require.NoError(t, err)
				assert.Equal(t, "multi", u.NickName)
				assert.Equal(t, "multi@example.com", u.Email)
				assert.Equal(t, 2, u.Type)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("users"))
				s := New(db)
				err := s.UpdateUser(context.Background(), tt.update)
				if tt.wantErr {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				if tt.check != nil {
					tt.check(t, s, tt.id)
				}
			})
		})
	}
}

func TestUser_GetUserByIds(t *testing.T) {
	tests := []struct {
		name    string
		ids     []int
		wantLen int
	}{
		{name: "多个ID", ids: []int{1, 2, 3}, wantLen: 3},
		{name: "单个ID", ids: []int{1}, wantLen: 1},
		{name: "包含不存在ID", ids: []int{1, 999}, wantLen: 1},
		{name: "空ID列表", ids: []int{}, wantLen: 0},
		{name: "全部不存在", ids: []int{999, 1000}, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("users"))
				s := New(db)
				ret, err := s.GetUserByIds(context.Background(), tt.ids)
				require.NoError(t, err)
				if tt.wantLen == 0 {
					assert.Empty(t, ret)
					return
				}
				assert.Len(t, ret, tt.wantLen)
				// 验证 map 的 key 正确
				for id := range ret {
					assert.Contains(t, tt.ids, id)
				}
			})
		})
	}
}

func TestUser_GetUserIDByOpenid(t *testing.T) {
	// fixture 中 openid 未设置，默认为空字符串，因此空 openid 会匹配到第一条记录
	tests := []struct {
		name    string
		openid  string
		wantErr bool
		wantId  int
	}{
		{name: "不存在的openid", openid: "nonexistent_openid", wantErr: true, wantId: 0},
		{name: "空openid匹配默认值", openid: "", wantErr: false, wantId: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("users"))
				s := New(db)
				uid, err := s.GetUserIDByOpenid(context.Background(), tt.openid)
				if tt.wantErr {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				assert.Equal(t, tt.wantId, uid)
			})
		})
	}
}
