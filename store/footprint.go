package store

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"app/store/model"
)

func (s *Store) GetFootprint(ctx context.Context, id int) (*model.Footprint, error) {
	query := "SELECT id, name, description, longitude, latitude, date, marker_color, categories, url, url_label, photos, created_at, updated_at FROM footprints WHERE id = ?"
	var m model.Footprint
	var catRaw, photoRaw []byte
	err := s.db.QueryRowContext(ctx, query, id).Scan(&m.Id, &m.Name, &m.Description, &m.Longitude, &m.Latitude, &m.Date, &m.MarkerColor, &catRaw, &m.Url, &m.UrlLabel, &photoRaw, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if err := m.ScanCategories(catRaw); err != nil {
		return nil, err
	}
	if err := m.ScanPhotos(photoRaw); err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *Store) ListFootprint(ctx context.Context, start int, num int) ([]*model.Footprint, error) {
	offset := (start - 1) * num
	rows, err := s.db.QueryContext(ctx, "SELECT id, name, description, longitude, latitude, date, marker_color, categories, url, url_label, photos, created_at, updated_at FROM footprints ORDER BY id DESC LIMIT ? OFFSET ?", num, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret := make([]*model.Footprint, 0)
	for rows.Next() {
		var item model.Footprint
		var catRaw, photoRaw []byte
		if err := rows.Scan(&item.Id, &item.Name, &item.Description, &item.Longitude, &item.Latitude, &item.Date, &item.MarkerColor, &catRaw, &item.Url, &item.UrlLabel, &photoRaw, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		_ = item.ScanCategories(catRaw)
		_ = item.ScanPhotos(photoRaw)
		tmp := item
		ret = append(ret, &tmp)
	}
	return ret, nil
}

func (s *Store) CountFootprintTotal(ctx context.Context) (int, error) {
	var total int
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM footprints").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (s *Store) CreateFootprint(ctx context.Context, md *model.Footprint) (int64, error) {
	catJSON, _ := json.Marshal(md.Categories)
	photoJSON, _ := json.Marshal(md.Photos)
	res, err := s.db.ExecContext(ctx,
		"INSERT INTO footprints (name, description, longitude, latitude, date, marker_color, categories, url, url_label, photos, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		md.Name, md.Description, md.Longitude, md.Latitude, md.Date, md.MarkerColor, catJSON, md.Url, md.UrlLabel, photoJSON, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) UpdateFootprint(ctx context.Context, md *model.UpdateFootprint) error {
	set := make([]string, 0)
	args := make([]any, 0)
	if v := md.Name; v != nil {
		set, args = append(set, "`name` = ?"), append(args, *v)
	}
	if v := md.Description; v != nil {
		set, args = append(set, "`description` = ?"), append(args, *v)
	}
	if v := md.Longitude; v != nil {
		set, args = append(set, "`longitude` = ?"), append(args, *v)
	}
	if v := md.Latitude; v != nil {
		set, args = append(set, "`latitude` = ?"), append(args, *v)
	}
	if v := md.Date; v != nil {
		set, args = append(set, "`date` = ?"), append(args, *v)
	}
	if v := md.MarkerColor; v != nil {
		set, args = append(set, "`marker_color` = ?"), append(args, *v)
	}
	if md.Categories != nil {
		catJSON, _ := json.Marshal(md.Categories)
		set, args = append(set, "`categories` = ?"), append(args, catJSON)
	}
	if v := md.Url; v != nil {
		set, args = append(set, "`url` = ?"), append(args, *v)
	}
	if v := md.UrlLabel; v != nil {
		set, args = append(set, "`url_label` = ?"), append(args, *v)
	}
	if md.PhotoUrls != nil {
		photos := model.PhotosFromURLs(md.PhotoUrls)
		photoJSON, _ := json.Marshal(photos)
		set, args = append(set, "`photos` = ?"), append(args, photoJSON)
	}
	if len(set) == 0 {
		return nil
	}
	args = append(args, md.Id)
	query := "UPDATE footprints SET " + strings.Join(set, ", ") + " WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func (s *Store) DeleteFootprint(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM footprints WHERE id = ?", id)
	return err
}

func (s *Store) ListAllFootprints(ctx context.Context) ([]*model.Footprint, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, name, description, longitude, latitude, date, marker_color, categories, url, url_label, photos, created_at, updated_at FROM footprints ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret := make([]*model.Footprint, 0)
	for rows.Next() {
		var item model.Footprint
		var catRaw, photoRaw []byte
		if err := rows.Scan(&item.Id, &item.Name, &item.Description, &item.Longitude, &item.Latitude, &item.Date, &item.MarkerColor, &catRaw, &item.Url, &item.UrlLabel, &photoRaw, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		_ = item.ScanCategories(catRaw)
		_ = item.ScanPhotos(photoRaw)
		tmp := item
		ret = append(ret, &tmp)
	}
	return ret, nil
}
