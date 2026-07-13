package store

import (
	"context"
	"database/sql"
	"math"
	"strconv"

	"app/store/model"
)

func (s *Store) GetRegion(ctx context.Context, regionId int) (*model.Region, error) {
	query := "SELECT region_id, parent_id, level, region_name, longitude, latitude, pinyin, az_no FROM regions WHERE region_id = ?"
	var m model.Region
	err := s.db.QueryRowContext(ctx, query, regionId).Scan(&m.RegionId, &m.ParentId, &m.Level, &m.RegionName, &m.Longitude, &m.Latitude, &m.Pinyin, &m.AzNo)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *Store) GetRegionByIds(ctx context.Context, ids []int) (map[int]model.Region, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	ph, args := In(ids)
	query := "SELECT region_id, parent_id, level, region_name, longitude, latitude, pinyin, az_no FROM regions WHERE region_id IN(" + ph + ")"
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rm := make(map[int]model.Region, len(ids))
	for rows.Next() {
		var m model.Region
		if err := rows.Scan(&m.RegionId, &m.ParentId, &m.Level, &m.RegionName, &m.Longitude, &m.Latitude, &m.Pinyin, &m.AzNo); err != nil {
			return nil, err
		}
		rm[m.RegionId] = m
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return rm, nil
}

func (s *Store) ListRegionByParent(ctx context.Context, parentId int) ([]*model.Region, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT region_id, parent_id, level, region_name, longitude, latitude, pinyin, az_no FROM regions WHERE parent_id = ? ORDER BY region_id", parentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret := make([]*model.Region, 0)
	for rows.Next() {
		var item model.Region
		if err := rows.Scan(&item.RegionId, &item.ParentId, &item.Level, &item.RegionName, &item.Longitude, &item.Latitude, &item.Pinyin, &item.AzNo); err != nil {
			return nil, err
		}
		tmp := item
		ret = append(ret, &tmp)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *Store) FindNearestCity(ctx context.Context, latitude, longitude float64) (*model.Region, *model.Region, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT region_id, parent_id, level, region_name, longitude, latitude, pinyin, az_no FROM regions WHERE level = 2")
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var best *model.Region
	bestDist := math.MaxFloat64
	for rows.Next() {
		var item model.Region
		if err := rows.Scan(&item.RegionId, &item.ParentId, &item.Level, &item.RegionName, &item.Longitude, &item.Latitude, &item.Pinyin, &item.AzNo); err != nil {
			return nil, nil, err
		}
		cityLng, err1 := strconv.ParseFloat(item.Longitude, 64)
		cityLat, err2 := strconv.ParseFloat(item.Latitude, 64)
		if err1 != nil || err2 != nil {
			continue
		}
		d := haversine(latitude, longitude, cityLat, cityLng)
		if d < bestDist {
			tmp := item
			best = &tmp
			bestDist = d
		}
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
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
