package store

import (
	"context"
	"database/sql"
	"math"
	"strconv"

	"app/pkg/sqlext"
	"app/store/model"
)

// regionColumns regions 表查询列
const regionColumns = "region_id, parent_id, level, region_name, longitude, latitude, pinyin, az_no"

// GetRegion 按区域 ID 查询，不存在时返回 sql.ErrNoRows
func (s *Store) GetRegion(ctx context.Context, regionId int) (*model.Region, error) {
	region, err := sqlext.QueryRow[model.Region](ctx, s.db, "select "+regionColumns+" from regions where region_id = ?", regionId)
	if err != nil {
		return nil, err
	}
	return &region, nil
}

// GetRegionByIds 按区域 ID 列表查询，返回以区域 ID 为键的map
func (s *Store) GetRegionByIds(ctx context.Context, ids []int) (map[int]model.Region, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	q := sqlext.NewBuilder().Select(regionColumns).From("regions").Where("region_id in (?)", ids)
	list, err := sqlext.Query[model.Region](ctx, s.db, q.SQL(), q.Args()...)
	if err != nil {
		return nil, err
	}

	rm := make(map[int]model.Region, len(list))
	for _, region := range list {
		rm[region.RegionId] = region
	}
	return rm, nil
}

// ListRegionByParent 查询某上级区域下的所有区域，按区域 ID 正序
func (s *Store) ListRegionByParent(ctx context.Context, parentId int) ([]model.Region, error) {
	q := sqlext.NewBuilder().
		Select(regionColumns).
		From("regions").
		Where("parent_id = ?", parentId).
		OrderBy("region_id")

	return sqlext.Query[model.Region](ctx, s.db, q.SQL(), q.Args()...)
}

// FindNearestCity 按经纬度查找最近的市级区域及其所属省份，无数据时返回 sql.ErrNoRows
func (s *Store) FindNearestCity(ctx context.Context, latitude, longitude float64) (*model.Region, *model.Region, error) {
	q := sqlext.NewBuilder().Select(regionColumns).From("regions").Where("level = ?", 2)

	list, err := sqlext.Query[model.Region](ctx, s.db, q.SQL(), q.Args()...)
	if err != nil {
		return nil, nil, err
	}

	var best *model.Region
	bestDist := math.MaxFloat64
	for i := range list {
		item := &list[i]
		cityLng, err1 := strconv.ParseFloat(item.Longitude, 64)
		cityLat, err2 := strconv.ParseFloat(item.Latitude, 64)
		if err1 != nil || err2 != nil {
			continue
		}

		if d := haversine(latitude, longitude, cityLat, cityLng); d < bestDist {
			best = item
			bestDist = d
		}
	}
	if best == nil {
		return nil, nil, sql.ErrNoRows
	}

	province, err := s.GetRegion(ctx, best.ParentId)
	if err != nil {
		return nil, nil, err
	}
	return best, province, nil
}

// haversine 计算两个经纬度坐标之间的球面距离（公里）
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371.0
	rad := func(v float64) float64 { return v * math.Pi / 180 }
	dLat := rad(lat2 - lat1)
	dLon := rad(lon2 - lon1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(rad(lat1))*math.Cos(rad(lat2))*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}
