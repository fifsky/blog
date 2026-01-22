package store

import (
	"context"
	"strings"
	"time"

	"app/store/model"
)

func (s *Store) GetPhoto(ctx context.Context, id int) (*model.Photo, error) {
	query := "SELECT id, title, description, src, thumbnail, province, city, created_at, updated_at FROM photos WHERE id = ?"
	var m model.Photo
	err := s.db.QueryRowContext(ctx, query, id).Scan(&m.Id, &m.Title, &m.Description, &m.Src, &m.Thumbnail, &m.Province, &m.City, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *Store) ListPhoto(ctx context.Context, start int, num int) ([]*model.Photo, error) {
	offset := (start - 1) * num
	rows, err := s.db.QueryContext(ctx, "SELECT id, title, description, src, thumbnail, province, city, created_at, updated_at FROM photos ORDER BY id DESC LIMIT ? OFFSET ?", num, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret := make([]*model.Photo, 0)
	for rows.Next() {
		var item model.Photo
		if err := rows.Scan(&item.Id, &item.Title, &item.Description, &item.Src, &item.Thumbnail, &item.Province, &item.City, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		tmp := item
		ret = append(ret, &tmp)
	}
	return ret, nil
}

func (s *Store) CountPhotoTotal(ctx context.Context) (int, error) {
	var total int
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM photos").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (s *Store) CreatePhoto(ctx context.Context, md *model.Photo) (int64, error) {
	res, err := s.db.ExecContext(ctx, "INSERT INTO photos (title, description, src, thumbnail, province, city, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		md.Title, md.Description, md.Src, md.Thumbnail, md.Province, md.City, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) UpdatePhoto(ctx context.Context, md *model.UpdatePhoto) error {
	set := make([]string, 0)
	args := make([]any, 0)
	if v := md.Title; v != nil {
		set, args = append(set, "`title` = ?"), append(args, *v)
	}
	if v := md.Description; v != nil {
		set, args = append(set, "`description` = ?"), append(args, *v)
	}
	if v := md.Province; v != nil {
		set, args = append(set, "`province` = ?"), append(args, *v)
	}
	if v := md.City; v != nil {
		set, args = append(set, "`city` = ?"), append(args, *v)
	}
	if len(set) == 0 {
		return nil
	}
	args = append(args, md.Id)
	query := "UPDATE photos SET " + strings.Join(set, ", ") + " WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func (s *Store) DeletePhoto(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM photos WHERE id = ?", id)
	return err
}

func (s *Store) ListPhotoByCity(ctx context.Context, city string) ([]*model.Photo, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, title, description, src, thumbnail, province, city, created_at, updated_at FROM photos WHERE city = ? ORDER BY id DESC", city)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret := make([]*model.Photo, 0)
	for rows.Next() {
		var item model.Photo
		if err := rows.Scan(&item.Id, &item.Title, &item.Description, &item.Src, &item.Thumbnail, &item.Province, &item.City, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		tmp := item
		ret = append(ret, &tmp)
	}
	return ret, nil
}
