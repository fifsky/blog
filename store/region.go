package store

import (
	"context"

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
	return ret, nil
}

// ListProvincesWithPhotos returns provinces that have at least one photo
func (s *Store) ListProvincesWithPhotos(ctx context.Context) ([]*model.Region, error) {
	query := `
		SELECT DISTINCT r.region_id, r.parent_id, r.level, r.region_name, r.longitude, r.latitude, r.pinyin, r.az_no 
		FROM regions r 
		INNER JOIN photos p ON r.region_id = p.province 
		WHERE r.level = 1 
		ORDER BY r.region_id
	`
	rows, err := s.db.QueryContext(ctx, query)
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
	return ret, nil
}

// ListCitiesWithPhotos returns cities that have at least one photo
func (s *Store) ListCitiesWithPhotos(ctx context.Context) ([]*model.Region, error) {
	query := `
		SELECT DISTINCT r.region_id, r.parent_id, r.level, r.region_name, r.longitude, r.latitude, r.pinyin, r.az_no 
		FROM regions r 
		INNER JOIN photos p ON r.region_id = p.city 
		WHERE r.level = 2 
		ORDER BY r.region_id
	`
	rows, err := s.db.QueryContext(ctx, query)
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
	return ret, nil
}
