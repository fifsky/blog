package store

import (
	"context"
	"strings"

	"app/store/model"
)

func (s *Store) ListGuestbook(ctx context.Context, keyword string, start int, num int) ([]model.Guestbook, error) {
	offset := (start - 1) * num
	query := "select id,name,content,ip,top,created_at from guestbook"
	args := []interface{}{}

	if keyword != "" {
		query += " where (name like ? or content like ?)"
		like := "%" + keyword + "%"
		args = append(args, like, like)
	}

	query += " order by top desc, created_at desc limit ? offset ?"
	args = append(args, num, offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	gs := make([]model.Guestbook, 0)
	for rows.Next() {
		var g model.Guestbook
		if err := rows.Scan(&g.Id, &g.Name, &g.Content, &g.Ip, &g.Top, &g.CreatedAt); err != nil {
			return nil, err
		}
		gs = append(gs, g)
	}
	return gs, nil
}

func (s *Store) CountGuestbookTotal(ctx context.Context, keyword string) (int, error) {
	query := "select count(*) from guestbook"
	args := []interface{}{}

	if keyword != "" {
		query += " where (name like ? or content like ?)"
		like := "%" + keyword + "%"
		args = append(args, like, like)
	}

	var total int
	err := s.db.QueryRowContext(ctx, query, args...).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (s *Store) CreateGuestbook(ctx context.Context, g *model.Guestbook) (int64, error) {
	res, err := s.db.ExecContext(ctx, "insert into guestbook (name,content,ip,created_at) values (?,?,?,?)", g.Name, g.Content, g.Ip, g.CreatedAt)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) GetGuestbook(ctx context.Context, id int64) (*model.Guestbook, error) {
	query := "select id,name,content,ip,top,created_at from guestbook where id = ?"
	var g model.Guestbook
	err := s.db.QueryRowContext(ctx, query, id).Scan(&g.Id, &g.Name, &g.Content, &g.Ip, &g.Top, &g.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (s *Store) DeleteGuestbook(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return nil
	}
	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	query := "delete from guestbook where id in (" + strings.Join(placeholders, ",") + ")"
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}
