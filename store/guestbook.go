package store

import (
	"context"

	"app/store/model"
)

func (s *Store) ListGuestbook(ctx context.Context, start int, num int) ([]model.Guestbook, error) {
	offset := (start - 1) * num
	rows, err := s.db.QueryContext(ctx, "select id,name,content,ip,created_at from guestbook order by id desc limit ? offset ?", num, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	gs := make([]model.Guestbook, 0)
	for rows.Next() {
		var g model.Guestbook
		if err := rows.Scan(&g.Id, &g.Name, &g.Content, &g.Ip, &g.CreatedAt); err != nil {
			return nil, err
		}
		gs = append(gs, g)
	}
	return gs, nil
}

func (s *Store) CountGuestbookTotal(ctx context.Context) (int, error) {
	var total int
	err := s.db.QueryRowContext(ctx, "select count(*) from guestbook").Scan(&total)
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
