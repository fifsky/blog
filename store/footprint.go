package store

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"app/pkg/sqlext"
	"app/store/model"
)

// footprintColumns footprints 表查询列。categories/photos 是 JSON 列，
// 加别名后由 footprintRow 先扫描成原始字节再解析，避免直接落到 []string / []FootprintPhoto 上
const footprintColumns = "id, name, description, longitude, latitude, date, marker_color," +
	" categories as categories_json, url, url_label, photos as photos_json, created_at, updated_at"

// footprintRow 是 footprints 表的扫描载体：嵌入 Footprint 复用其字段匹配，
// JSON 列单独用 []byte 接收后解析进 Categories / Photos
type footprintRow struct {
	model.Footprint
	CategoriesJSON []byte `db:"categories_json"`
	PhotosJSON     []byte `db:"photos_json"`
}

// footprint 把 JSON 列解析到嵌入式 Footprint（ScanCategories/ScanPhotos 由它提升而来）
func (r *footprintRow) footprint() (*model.Footprint, error) {
	if err := r.ScanCategories(r.CategoriesJSON); err != nil {
		return nil, err
	}
	if err := r.ScanPhotos(r.PhotosJSON); err != nil {
		return nil, err
	}
	return &r.Footprint, nil
}

// GetFootprint 按 ID 查询足迹，不存在时返回 sql.ErrNoRows
func (s *Store) GetFootprint(ctx context.Context, id int) (*model.Footprint, error) {
	row, err := sqlext.QueryRow[footprintRow](ctx, s.db, "select "+footprintColumns+" from footprints where id = ?", id)
	if err != nil {
		return nil, err
	}
	return row.footprint()
}

// ListFootprint 分页查询足迹，按 ID 倒序
func (s *Store) ListFootprint(ctx context.Context, start int, num int) ([]*model.Footprint, error) {
	q := sqlext.NewBuilder().
		Select(footprintColumns).
		From("footprints").
		OrderBy("id desc").
		Limit(num).
		Offset((start - 1) * num)

	rows, err := sqlext.Query[footprintRow](ctx, s.db, q.SQL(), q.Args()...)
	if err != nil {
		return nil, err
	}
	return footprints(rows)
}

// ListAllFootprints 查询全部足迹，按 ID 倒序
func (s *Store) ListAllFootprints(ctx context.Context) ([]*model.Footprint, error) {
	q := sqlext.NewBuilder().Select(footprintColumns).From("footprints").OrderBy("id desc")

	rows, err := sqlext.Query[footprintRow](ctx, s.db, q.SQL(), q.Args()...)
	if err != nil {
		return nil, err
	}
	return footprints(rows)
}

// CountFootprintTotal 统计足迹总数
func (s *Store) CountFootprintTotal(ctx context.Context) (int, error) {
	q := sqlext.NewBuilder().Select("count(*)").From("footprints")
	return sqlext.QueryRow[int](ctx, s.db, q.SQL(), q.Args()...)
}

// CreateFootprint 新增足迹，返回自增 ID
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

// UpdateFootprint 按需更新足迹字段，只更新传入的非空字段
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

// DeleteFootprint 按 ID 删除足迹
func (s *Store) DeleteFootprint(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM footprints WHERE id = ?", id)
	return err
}

// footprints 把扫描出的行解析成足迹切片，空结果返回空切片而不是 nil
func footprints(rows []footprintRow) ([]*model.Footprint, error) {
	ret := make([]*model.Footprint, 0, len(rows))
	for i := range rows {
		item, err := rows[i].footprint()
		if err != nil {
			return nil, err
		}
		ret = append(ret, item)
	}
	return ret, nil
}
