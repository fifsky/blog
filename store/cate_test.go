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

func TestCate_GetAllCates(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("cates", "posts"))
		s := New(db)
		ret, err := s.GetAllCates(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, ret)
	})
}

func TestCate_GetCate(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		wantErr bool
		want    string
	}{
		{name: "存在分类", id: 1, wantErr: false, want: "默认分类"},
		{name: "不存在分类", id: 999, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("cates"))
				s := New(db)
				ret, err := s.GetCate(context.Background(), tt.id)
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

func TestCate_CreateCate(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("cates"))
		s := New(db)

		now := time.Now()
		c := &model.Cate{
			Name:      "技术",
			Desc:      "技术文章",
			Domain:    "tech",
			CreatedAt: now,
			UpdatedAt: now,
		}

		id, err := s.CreateCate(context.Background(), c)
		require.NoError(t, err)
		assert.Greater(t, id, int64(0))

		// 验证创建的数据
		got, err := s.GetCate(context.Background(), int(id))
		require.NoError(t, err)
		assert.Equal(t, "技术", got.Name)
		assert.Equal(t, "tech", got.Domain)
		assert.Equal(t, "技术文章", got.Desc)
	})

	t.Run("重复domain报错", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("cates"))
			s := New(db)

			now := time.Now()
			c := &model.Cate{
				Name:      "重复分类",
				Desc:      "desc",
				Domain:    "default", // 已存在的domain
				CreatedAt: now,
				UpdatedAt: now,
			}
			_, err := s.CreateCate(context.Background(), c)
			require.Error(t, err)
		})
	})
}

func TestCate_UpdateCate(t *testing.T) {
	tests := []struct {
		name   string
		update *model.UpdateCate
		check  func(t *testing.T, c *model.Cate)
	}{
		{
			name: "更新名称",
			update: &model.UpdateCate{
				Id:   1,
				Name: new("新名称"),
			},
			check: func(t *testing.T, c *model.Cate) {
				assert.Equal(t, "新名称", c.Name)
			},
		},
		{
			name: "更新描述",
			update: &model.UpdateCate{
				Id:   1,
				Desc: new("新描述"),
			},
			check: func(t *testing.T, c *model.Cate) {
				assert.Equal(t, "新描述", c.Desc)
			},
		},
		{
			name: "更新域名",
			update: &model.UpdateCate{
				Id:     1,
				Domain: new("newdomain"),
			},
			check: func(t *testing.T, c *model.Cate) {
				assert.Equal(t, "newdomain", c.Domain)
			},
		},
		{
			name: "更新多个字段",
			update: &model.UpdateCate{
				Id:     1,
				Name:   new("多字段"),
				Desc:   new("多描述"),
				Domain: new("multi"),
			},
			check: func(t *testing.T, c *model.Cate) {
				assert.Equal(t, "多字段", c.Name)
				assert.Equal(t, "多描述", c.Desc)
				assert.Equal(t, "multi", c.Domain)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("cates"))
				s := New(db)
				err := s.UpdateCate(context.Background(), tt.update)
				require.NoError(t, err)

				got, err := s.GetCate(context.Background(), tt.update.Id)
				require.NoError(t, err)
				tt.check(t, got)
			})
		})
	}
}

func TestCate_DeleteCate(t *testing.T) {
	t.Run("删除存在的分类", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("cates"))
			s := New(db)
			err := s.DeleteCate(context.Background(), 1)
			require.NoError(t, err)

			// 验证已删除
			_, err = s.GetCate(context.Background(), 1)
			require.Error(t, err)
		})
	})

	t.Run("删除不存在的分类不报错", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("cates"))
			s := New(db)
			err := s.DeleteCate(context.Background(), 999)
			require.NoError(t, err)
		})
	})
}

func TestCate_GetCatesByIds(t *testing.T) {
	tests := []struct {
		name      string
		ids       []int
		wantEmpty bool
	}{
		{name: "存在ID", ids: []int{1}, wantEmpty: false},
		{name: "包含不存在ID", ids: []int{1, 999}, wantEmpty: false},
		{name: "空ID列表", ids: []int{}, wantEmpty: true},
		{name: "全部不存在", ids: []int{999}, wantEmpty: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("cates"))
				s := New(db)
				ret, err := s.GetCatesByIds(context.Background(), tt.ids)
				require.NoError(t, err)
				if tt.wantEmpty {
					assert.Empty(t, ret)
					return
				}
				assert.NotEmpty(t, ret)
			})
		})
	}
}

func TestCate_PostsCount(t *testing.T) {
	tests := []struct {
		name      string
		cateId    int
		wantEmpty bool
	}{
		{name: "有文章的分类", cateId: 1, wantEmpty: false},
		{name: "无文章的分类", cateId: 999, wantEmpty: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("cates", "posts"))
				s := New(db)
				count, err := s.PostsCount(context.Background(), tt.cateId)
				require.NoError(t, err)
				if tt.wantEmpty {
					assert.Equal(t, 0, count)
				} else {
					assert.Greater(t, count, 0)
				}
			})
		})
	}
}
