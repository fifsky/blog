package store

import (
	"context"
	"strings"

	"app/store/model"
)

func (s *Store) GetAllLinks(ctx context.Context) ([]*model.Link, error) {
	rows, err := s.db.QueryContext(ctx, "select id,name,url,`desc`,status,created_at,updated_at from links order by id asc")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret := make([]*model.Link, 0)
	for rows.Next() {
		var item model.Link
		if err := rows.Scan(&item.Id, &item.Name, &item.Url, &item.Desc, &item.Status, &item.CreatedAt, &item.UpdatedAt); err != nil {
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

// GetApprovedLinks 获取审核通过的链接列表
func (s *Store) GetApprovedLinks(ctx context.Context) ([]*model.Link, error) {
	rows, err := s.db.QueryContext(ctx, "select id,name,url,`desc`,status,created_at,updated_at from links where status = ? order by id asc", model.LinkStatusApproved)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret := make([]*model.Link, 0)
	for rows.Next() {
		var item model.Link
		if err := rows.Scan(&item.Id, &item.Name, &item.Url, &item.Desc, &item.Status, &item.CreatedAt, &item.UpdatedAt); err != nil {
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

// GetLink 根据 ID 获取链接信息
func (s *Store) GetLink(ctx context.Context, id int) (*model.Link, error) {
	var item model.Link
	err := s.db.QueryRowContext(ctx, "select id,name,url,`desc`,status,created_at,updated_at from links where id = ?", id).
		Scan(&item.Id, &item.Name, &item.Url, &item.Desc, &item.Status, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *Store) CreateLink(ctx context.Context, link *model.Link) (int64, error) {
	res, err := s.db.ExecContext(ctx, "insert into links (name,url,`desc`,status,created_at,updated_at) values (?,?,?,?,?,?)",
		link.Name, link.Url, link.Desc, link.Status, link.CreatedAt, link.UpdatedAt)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) UpdateLink(ctx context.Context, link *model.UpdateLink) error {
	set := make([]string, 0)
	args := make([]any, 0)
	if v := link.Name; v != nil {
		set, args = append(set, "`name` = ?"), append(args, *v)
	}
	if v := link.Url; v != nil {
		set, args = append(set, "`url` = ?"), append(args, *v)
	}
	if v := link.Desc; v != nil {
		set, args = append(set, "`desc` = ?"), append(args, *v)
	}
	if v := link.Status; v != nil {
		set, args = append(set, "`status` = ?"), append(args, *v)
	}
	args = append(args, link.Id)
	query := "update links set " + strings.Join(set, ", ") + " where id = ?"
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func (s *Store) DeleteLink(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, "delete from links where id = ?", id)
	return err
}
