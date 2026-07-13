package store

import (
	"context"
	"testing"
	"time"

	"app/pkg/dbunit"
	"app/store/model"
	"app/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLink_GetAllLinks(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("links"))
		s := New(db)
		ret, err := s.GetAllLinks(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		assert.NotEmpty(t, ret)
	})
}

func TestLink_GetApprovedLinks(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("links"))
		s := New(db)
		ret, err := s.GetApprovedLinks(context.Background())
		require.NoError(t, err)
		assert.NotEmpty(t, ret)
		for _, link := range ret {
			assert.Equal(t, model.LinkStatusApproved, link.Status)
		}
	})

	t.Run("只有APPROVED状态", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), nil)
			s := New(db)

			// 创建一条 PENDING 状态的链接
			now := time.Now()
			_, err := s.CreateLink(context.Background(), &model.Link{
				Name:      "pending_link",
				Url:       "http://pending.com",
				Desc:      "pending",
				Status:    model.LinkStatusPending,
				CreatedAt: now,
				UpdatedAt: now,
			})
			require.NoError(t, err)

			// 创建一条 APPROVED 状态的链接
			_, err = s.CreateLink(context.Background(), &model.Link{
				Name:      "approved_link",
				Url:       "http://approved.com",
				Desc:      "approved",
				Status:    model.LinkStatusApproved,
				CreatedAt: now,
				UpdatedAt: now,
			})
			require.NoError(t, err)

			ret, err := s.GetApprovedLinks(context.Background())
			require.NoError(t, err)
			assert.Len(t, ret, 1)
			assert.Equal(t, "approved_link", ret[0].Name)
		})
	})
}

func TestLink_GetLink(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		wantErr bool
		want    string
	}{
		{name: "存在链接", id: 1, wantErr: false, want: "圆子"},
		{name: "不存在链接", id: 999, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("links"))
				s := New(db)
				ret, err := s.GetLink(context.Background(), tt.id)
				if tt.wantErr {
					require.Error(t, err)
					assert.Nil(t, ret)
					return
				}
				require.NoError(t, err)
				assert.NotNil(t, ret)
				assert.Equal(t, tt.want, ret.Name)
			})
		})
	}
}

func TestLink_CreateLink(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), nil)
		s := New(db)

		now := time.Now()
		link := &model.Link{
			Name:      "新链接",
			Url:       "https://example.com",
			Desc:      "测试链接",
			Status:    model.LinkStatusPending,
			CreatedAt: now,
			UpdatedAt: now,
		}

		id, err := s.CreateLink(context.Background(), link)
		require.NoError(t, err)
		assert.Greater(t, id, int64(0))

		// 验证创建的数据
		got, err := s.GetLink(context.Background(), int(id))
		require.NoError(t, err)
		assert.Equal(t, "新链接", got.Name)
		assert.Equal(t, "https://example.com", got.Url)
		assert.Equal(t, "测试链接", got.Desc)
		assert.Equal(t, model.LinkStatusPending, got.Status)
	})
}

func TestLink_UpdateLink(t *testing.T) {
	tests := []struct {
		name   string
		update *model.UpdateLink
		check  func(t *testing.T, l *model.Link)
	}{
		{
			name: "更新名称",
			update: &model.UpdateLink{
				Id:   1,
				Name: new("新名称"),
			},
			check: func(t *testing.T, l *model.Link) {
				assert.Equal(t, "新名称", l.Name)
			},
		},
		{
			name: "更新URL",
			update: &model.UpdateLink{
				Id:  1,
				Url: new("https://newurl.com"),
			},
			check: func(t *testing.T, l *model.Link) {
				assert.Equal(t, "https://newurl.com", l.Url)
			},
		},
		{
			name: "更新描述",
			update: &model.UpdateLink{
				Id:   1,
				Desc: new("新描述"),
			},
			check: func(t *testing.T, l *model.Link) {
				assert.Equal(t, "新描述", l.Desc)
			},
		},
		{
			name: "更新状态为APPROVED",
			update: &model.UpdateLink{
				Id:     1,
				Status: new(model.LinkStatusApproved),
			},
			check: func(t *testing.T, l *model.Link) {
				assert.Equal(t, model.LinkStatusApproved, l.Status)
			},
		},
		{
			name: "更新多个字段",
			update: &model.UpdateLink{
				Id:     1,
				Name:   new("多字段"),
				Url:    new("https://multi.com"),
				Desc:   new("多描述"),
				Status: new(model.LinkStatusPending),
			},
			check: func(t *testing.T, l *model.Link) {
				assert.Equal(t, "多字段", l.Name)
				assert.Equal(t, "https://multi.com", l.Url)
				assert.Equal(t, "多描述", l.Desc)
				assert.Equal(t, model.LinkStatusPending, l.Status)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("links"))
				s := New(db)
				err := s.UpdateLink(context.Background(), tt.update)
				require.NoError(t, err)

				got, err := s.GetLink(context.Background(), tt.update.Id)
				require.NoError(t, err)
				tt.check(t, got)
			})
		})
	}
}

func TestLink_DeleteLink(t *testing.T) {
	t.Run("删除存在链接", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("links"))
			s := New(db)
			err := s.DeleteLink(context.Background(), 1)
			require.NoError(t, err)

			_, err = s.GetLink(context.Background(), 1)
			require.Error(t, err)
		})
	})

	t.Run("删除不存在链接不报错", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("links"))
			s := New(db)
			err := s.DeleteLink(context.Background(), 999)
			require.NoError(t, err)
		})
	})
}
