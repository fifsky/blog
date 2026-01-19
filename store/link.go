package store

import (
	"context"
	"strings"

	"app/store/model"
)

func (s *Store) GetAllLinks(ctx context.Context) ([]*model.Link, error) {
	rows, err := s.db.QueryContext(ctx, "select id,name,url,`desc`,created_at from links")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret := make([]*model.Link, 0)
	for rows.Next() {
		var item model.Link
		if err := rows.Scan(&item.Id, &item.Name, &item.Url, &item.Desc, &item.CreatedAt); err != nil {
			return nil, err
		}
		tmp := item
		ret = append(ret, &tmp)
	}
	return ret, nil
}

func (s *Store) CreateLink(ctx context.Context, link *model.Link) (int64, error) {
	res, err := s.db.ExecContext(ctx, "insert into links (name,url,`desc`,created_at) values (?,?,?,?)",
		link.Name, link.Url, link.Desc, link.CreatedAt)
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
	args = append(args, link.Id)
	query := "update links set " + strings.Join(set, ", ") + " where id = ?"
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func (s *Store) DeleteLink(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, "delete from links where id = ?", id)
	return err
}
