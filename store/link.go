package store

import (
	"context"
	"fmt"
	"strings"

	"app/store/model"
)

func (s *Store) GetAllLinks(ctx context.Context) ([]*model.Link, error) {
	rows, err := s.db.QueryContext(ctx, `select id,name,url,"desc",created_at from blog.links`)
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
	var id int64
	err := s.db.QueryRowContext(ctx, `insert into blog.links (name,url,"desc",created_at) values ($1,$2,$3,$4) RETURNING id`,
		link.Name, link.Url, link.Desc, link.CreatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Store) UpdateLink(ctx context.Context, link *model.UpdateLink) error {
	set := make([]string, 0)
	args := make([]any, 0)

	if v := link.Name; v != nil {
		set, args = append(set, "name = "+placeholder(len(args)+1)), append(args, *v)
	}
	if v := link.Url; v != nil {
		set, args = append(set, "url = "+placeholder(len(args)+1)), append(args, *v)
	}
	if v := link.Desc; v != nil {
		set, args = append(set, `"desc" = `+placeholder(len(args)+1)), append(args, *v)
	}
	args = append(args, link.Id)
	query := fmt.Sprintf("update blog.links set %s where id = %s", strings.Join(set, ", "), placeholder(len(args)))
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func (s *Store) DeleteLink(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, "delete from blog.links where id = $1", id)
	return err
}
