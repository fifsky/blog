package store

import (
	"context"
	"testing"

	"app/pkg/dbunit"
	"app/store/model"
	"app/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFootprint_GetFootprint(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		wantErr bool
		want    string
	}{
		{name: "存在足迹", id: 1, wantErr: false, want: "上海"},
		{name: "存在足迹2", id: 2, wantErr: false, want: "杭州"},
		{name: "不存在足迹", id: 999, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("footprints"))
				s := New(db)
				ret, err := s.GetFootprint(context.Background(), tt.id)
				if tt.wantErr {
					require.Error(t, err)
					assert.Nil(t, ret)
					return
				}
				require.NoError(t, err)
				assert.NotNil(t, ret)
				assert.Equal(t, tt.want, ret.Name)
				// 验证 categories 和 photos 正确解析
				assert.NotEmpty(t, ret.Categories)
				assert.NotEmpty(t, ret.Photos)
			})
		})
	}
}

func TestFootprint_ListFootprint(t *testing.T) {
	tests := []struct {
		name    string
		start   int
		num     int
		wantLen int
	}{
		{name: "全部第一页", start: 1, num: 10, wantLen: 2},
		{name: "每页1条第1页", start: 1, num: 1, wantLen: 1},
		{name: "每页1条第3页无数据", start: 3, num: 1, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("footprints"))
				s := New(db)
				ret, err := s.ListFootprint(context.Background(), tt.start, tt.num)
				require.NoError(t, err)
				assert.Len(t, ret, tt.wantLen)
				// 验证按 id DESC 排序
				if len(ret) > 1 {
					assert.Greater(t, ret[0].Id, ret[1].Id)
				}
			})
		})
	}
}

func TestFootprint_CountFootprintTotal(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("footprints"))
		s := New(db)
		total, err := s.CountFootprintTotal(context.Background())
		require.NoError(t, err)
		assert.Equal(t, 2, total)
	})

	t.Run("空表", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), nil)
			s := New(db)
			total, err := s.CountFootprintTotal(context.Background())
			require.NoError(t, err)
			assert.Equal(t, 0, total)
		})
	})
}

func TestFootprint_CreateFootprint(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), nil)
		s := New(db)

		beforeTotal, err := s.CountFootprintTotal(context.Background())
		require.NoError(t, err)

		md := &model.Footprint{
			Name:        "北京",
			Description: "首都",
			Longitude:   "116.407396",
			Latitude:    "39.904200",
			Date:        "2024-08-01",
			MarkerColor: "red",
			Categories:  []string{"想去", "城市"},
			Url:         "https://example.com",
			UrlLabel:    "详情",
			Photos:      model.PhotosFromURLs([]string{"https://example.com/beijing.jpg"}),
		}

		id, err := s.CreateFootprint(context.Background(), md)
		require.NoError(t, err)
		assert.Greater(t, id, int64(0))

		afterTotal, err := s.CountFootprintTotal(context.Background())
		require.NoError(t, err)
		assert.Equal(t, beforeTotal+1, afterTotal)

		// 验证创建的数据
		got, err := s.GetFootprint(context.Background(), int(id))
		require.NoError(t, err)
		assert.Equal(t, "北京", got.Name)
		assert.Equal(t, "首都", got.Description)
		assert.Equal(t, "116.407396", got.Longitude)
		assert.Equal(t, "39.904200", got.Latitude)
		assert.Equal(t, "2024-08-01", got.Date)
		assert.Equal(t, "red", got.MarkerColor)
		assert.Contains(t, got.Categories, "想去")
		assert.Contains(t, got.Categories, "城市")
		assert.Equal(t, "https://example.com", got.Url)
		assert.Equal(t, "详情", got.UrlLabel)
		assert.Len(t, got.Photos, 1)
		assert.Equal(t, "https://example.com/beijing.jpg", got.Photos[0].Src)
	})
}

func TestFootprint_UpdateFootprint(t *testing.T) {
	tests := []struct {
		name   string
		update *model.UpdateFootprint
		check  func(t *testing.T, f *model.Footprint)
	}{
		{
			name: "更新名称",
			update: &model.UpdateFootprint{
				Id:   1,
				Name: new("新名称"),
			},
			check: func(t *testing.T, f *model.Footprint) {
				assert.Equal(t, "新名称", f.Name)
			},
		},
		{
			name: "更新描述",
			update: &model.UpdateFootprint{
				Id:          1,
				Description: new("新描述"),
			},
			check: func(t *testing.T, f *model.Footprint) {
				assert.Equal(t, "新描述", f.Description)
			},
		},
		{
			name: "更新经纬度",
			update: &model.UpdateFootprint{
				Id:        1,
				Longitude: new("116.0"),
				Latitude:  new("40.0"),
			},
			check: func(t *testing.T, f *model.Footprint) {
				assert.Equal(t, "116.0", f.Longitude)
				assert.Equal(t, "40.0", f.Latitude)
			},
		},
		{
			name: "更新日期和颜色",
			update: &model.UpdateFootprint{
				Id:          1,
				Date:        new("2025-01-01"),
				MarkerColor: new("blue"),
			},
			check: func(t *testing.T, f *model.Footprint) {
				assert.Equal(t, "2025-01-01", f.Date)
				assert.Equal(t, "blue", f.MarkerColor)
			},
		},
		{
			name: "更新分类",
			update: &model.UpdateFootprint{
				Id:         1,
				Categories: []string{"新分类1", "新分类2"},
			},
			check: func(t *testing.T, f *model.Footprint) {
				assert.Contains(t, f.Categories, "新分类1")
				assert.Contains(t, f.Categories, "新分类2")
			},
		},
		{
			name: "更新URL",
			update: &model.UpdateFootprint{
				Id:       1,
				Url:      new("https://new.com"),
				UrlLabel: new("新链接"),
			},
			check: func(t *testing.T, f *model.Footprint) {
				assert.Equal(t, "https://new.com", f.Url)
				assert.Equal(t, "新链接", f.UrlLabel)
			},
		},
		{
			name: "更新照片",
			update: &model.UpdateFootprint{
				Id:        1,
				PhotoUrls: []string{"https://example.com/new1.jpg", "https://example.com/new2.jpg"},
			},
			check: func(t *testing.T, f *model.Footprint) {
				assert.Len(t, f.Photos, 2)
				assert.Equal(t, "https://example.com/new1.jpg", f.Photos[0].Src)
				assert.Equal(t, "https://example.com/new2.jpg", f.Photos[1].Src)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("footprints"))
				s := New(db)
				err := s.UpdateFootprint(context.Background(), tt.update)
				require.NoError(t, err)

				got, err := s.GetFootprint(context.Background(), tt.update.Id)
				require.NoError(t, err)
				tt.check(t, got)
			})
		})
	}

	t.Run("空更新不报错", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("footprints"))
			s := New(db)
			err := s.UpdateFootprint(context.Background(), &model.UpdateFootprint{Id: 1})
			require.NoError(t, err)
		})
	})
}

func TestFootprint_DeleteFootprint(t *testing.T) {
	t.Run("删除存在足迹", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("footprints"))
			s := New(db)
			err := s.DeleteFootprint(context.Background(), 1)
			require.NoError(t, err)

			_, err = s.GetFootprint(context.Background(), 1)
			require.Error(t, err)
		})
	})

	t.Run("删除不存在足迹不报错", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("footprints"))
			s := New(db)
			err := s.DeleteFootprint(context.Background(), 999)
			require.NoError(t, err)
		})
	})
}

func TestFootprint_ListAllFootprints(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("footprints"))
		s := New(db)
		ret, err := s.ListAllFootprints(context.Background())
		require.NoError(t, err)
		assert.Len(t, ret, 2)
		// 验证按 id DESC 排序
		assert.Greater(t, ret[0].Id, ret[1].Id)
		// 验证 categories 和 photos 正确解析
		for _, f := range ret {
			assert.NotEmpty(t, f.Categories)
			assert.NotEmpty(t, f.Photos)
		}
	})

	t.Run("空表", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), nil)
			s := New(db)
			ret, err := s.ListAllFootprints(context.Background())
			require.NoError(t, err)
			assert.Empty(t, ret)
		})
	})
}
