package store

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"app/pkg/dbunit"
	"app/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegion_GetRegion(t *testing.T) {
	tests := []struct {
		name     string
		regionId int
		wantErr  bool
		wantName string
		wantLvl  int
	}{
		{name: "省份", regionId: 310000, wantErr: false, wantName: "上海市", wantLvl: 1},
		{name: "城市", regionId: 310100, wantErr: false, wantName: "上海市", wantLvl: 2},
		{name: "不存在", regionId: 999, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("regions"))
				s := New(db)
				ret, err := s.GetRegion(context.Background(), tt.regionId)
				if tt.wantErr {
					require.Error(t, err)
					assert.Nil(t, ret)
					return
				}
				require.NoError(t, err)
				assert.NotNil(t, ret)
				assert.Equal(t, tt.wantName, ret.RegionName)
				assert.Equal(t, tt.wantLvl, ret.Level)
			})
		})
	}
}

func TestRegion_GetRegionByIds(t *testing.T) {
	tests := []struct {
		name      string
		ids       []int
		wantEmpty bool
	}{
		{name: "多个ID", ids: []int{310000, 330000, 310100}, wantEmpty: false},
		{name: "单个ID", ids: []int{310000}, wantEmpty: false},
		{name: "包含不存在ID", ids: []int{310000, 999}, wantEmpty: false},
		{name: "空ID列表", ids: []int{}, wantEmpty: true},
		{name: "全部不存在", ids: []int{999, 1000}, wantEmpty: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("regions"))
				s := New(db)
				ret, err := s.GetRegionByIds(context.Background(), tt.ids)
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

func TestRegion_ListRegionByParent(t *testing.T) {
	tests := []struct {
		name      string
		parentId  int
		wantEmpty bool
		wantLvl   int
	}{
		{name: "上海市下属", parentId: 310000, wantEmpty: false, wantLvl: 2},
		{name: "浙江省下属", parentId: 330000, wantEmpty: false, wantLvl: 2},
		{name: "无下属", parentId: 310100, wantEmpty: true},
		{name: "不存在父级", parentId: 999, wantEmpty: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("regions"))
				s := New(db)
				ret, err := s.ListRegionByParent(context.Background(), tt.parentId)
				require.NoError(t, err)
				if tt.wantEmpty {
					assert.Empty(t, ret)
					return
				}
				assert.NotEmpty(t, ret)
				for _, r := range ret {
					assert.Equal(t, tt.wantLvl, r.Level)
					assert.Equal(t, tt.parentId, r.ParentId)
				}
			})
		})
	}
}

func TestRegion_FindNearestCity(t *testing.T) {
	t.Run("找到最近城市-上海市中心坐标", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("regions"))
			s := New(db)
			// 使用上海市坐标（31.230416, 121.473701），应匹配上海市(310100)
			city, province, err := s.FindNearestCity(context.Background(), 31.230416, 121.473701)
			require.NoError(t, err)
			assert.NotNil(t, city)
			assert.NotNil(t, province)
			assert.Equal(t, 310100, city.RegionId)
			assert.Equal(t, "上海市", city.RegionName)
			assert.Equal(t, 310000, province.RegionId)
			assert.Equal(t, 1, province.Level)
		})
	})

	t.Run("找到最近城市-杭州坐标", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("regions"))
			s := New(db)
			// 使用杭州市坐标（30.274084, 120.15507），应匹配杭州市(330100)
			city, province, err := s.FindNearestCity(context.Background(), 30.274084, 120.15507)
			require.NoError(t, err)
			assert.NotNil(t, city)
			assert.NotNil(t, province)
			assert.Equal(t, 330100, city.RegionId)
			assert.Equal(t, "杭州市", city.RegionName)
			assert.Equal(t, 330000, province.RegionId)
		})
	})

	t.Run("找到最近城市-宁波附近坐标", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("regions"))
			s := New(db)
			// 使用宁波市坐标（29.868388, 121.549792），应匹配宁波市(330200)
			city, province, err := s.FindNearestCity(context.Background(), 29.868388, 121.549792)
			require.NoError(t, err)
			assert.NotNil(t, city)
			assert.NotNil(t, province)
			assert.Equal(t, 330200, city.RegionId)
			assert.Equal(t, "宁波市", city.RegionName)
			assert.Equal(t, 330000, province.RegionId)
		})
	})

	t.Run("无城市数据返回ErrNoRows", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), nil)
			s := New(db)
			city, province, err := s.FindNearestCity(context.Background(), 31.0, 121.0)
			require.Error(t, err)
			assert.True(t, errors.Is(err, sql.ErrNoRows))
			assert.Nil(t, city)
			assert.Nil(t, province)
		})
	})
}

func TestHaversine(t *testing.T) {
	tests := []struct {
		name       string
		lat1, lon1 float64
		lat2, lon2 float64
		wantApprox float64
		tolerance  float64
	}{
		{
			name:       "相同点距离为0",
			lat1:       31.230416,
			lon1:       121.473701,
			lat2:       31.230416,
			lon2:       121.473701,
			wantApprox: 0,
			tolerance:  0.1,
		},
		{
			name:       "上海到杭州约170km",
			lat1:       31.230416,
			lon1:       121.473701,
			lat2:       30.274084,
			lon2:       120.15507,
			wantApprox: 170,
			tolerance:  20,
		},
		{
			name:       "上海到宁波约150km",
			lat1:       31.230416,
			lon1:       121.473701,
			lat2:       29.868388,
			lon2:       121.549792,
			wantApprox: 150,
			tolerance:  10,
		},
		{
			name:       "赤道上1度约111km",
			lat1:       0,
			lon1:       0,
			lat2:       0,
			lon2:       1,
			wantApprox: 111,
			tolerance:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dist := haversine(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			assert.InDelta(t, tt.wantApprox, dist, tt.tolerance)
		})
	}
}
