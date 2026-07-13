package store

import (
	"context"
	"strings"
	"time"

	"app/store/model"
)

func (s *Store) GetRemind(ctx context.Context, id int) (*model.Remind, error) {
	query := "select id,cron,content,status,next_time,created_at,updated_at from reminds where id = ?"
	var m model.Remind
	err := s.db.QueryRowContext(ctx, query, id).Scan(&m.Id, &m.Cron, &m.Content, &m.Status, &m.NextTime, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *Store) ListRemind(ctx context.Context, start int, num int) ([]*model.Remind, error) {
	offset := (start - 1) * num
	rows, err := s.db.QueryContext(ctx, "select id,cron,content,status,next_time,created_at,updated_at from reminds order by id desc limit ? offset ?", num, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret := make([]*model.Remind, 0)
	for rows.Next() {
		var item model.Remind
		if err := rows.Scan(&item.Id, &item.Cron, &item.Content, &item.Status, &item.NextTime, &item.CreatedAt, &item.UpdatedAt); err != nil {
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

func (s *Store) RemindAll(ctx context.Context) ([]model.Remind, error) {
	rows, err := s.db.QueryContext(ctx, "select id,cron,content,status,next_time,created_at,updated_at from reminds where status in (?, ?) order by id desc", model.RemindStatusActive, model.RemindStatusPending)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret := make([]model.Remind, 0)
	for rows.Next() {
		var item model.Remind
		if err := rows.Scan(&item.Id, &item.Cron, &item.Content, &item.Status, &item.NextTime, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		ret = append(ret, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *Store) CountRemindTotal(ctx context.Context) (int, error) {
	var total int
	err := s.db.QueryRowContext(ctx, "select count(*) from reminds").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (s *Store) UpdateRemindStatus(ctx context.Context, id int, status model.RemindStatus) error {
	_, err := s.db.ExecContext(ctx, "update reminds set status = ? where id = ?", status, id)
	return err
}

func (s *Store) UpdateRemindNextTime(ctx context.Context, id int, nextTime time.Time) error {
	_, err := s.db.ExecContext(ctx, "update reminds set next_time = ? where id = ?", nextTime, id)
	return err
}

func (s *Store) CreateRemind(ctx context.Context, md *model.Remind) (int64, error) {
	res, err := s.db.ExecContext(ctx, "insert into reminds (cron,content,status,next_time,created_at,updated_at) values (?,?,?,?,?,?)",
		md.Cron, md.Content, md.Status, md.NextTime, md.CreatedAt, md.UpdatedAt)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) UpdateRemind(ctx context.Context, md *model.UpdateRemind) error {
	set := make([]string, 0)
	args := make([]any, 0)
	if v := md.Cron; v != nil {
		set, args = append(set, "`cron` = ?"), append(args, *v)
	}
	if v := md.Content; v != nil {
		set, args = append(set, "`content` = ?"), append(args, *v)
	}
	if v := md.Status; v != nil {
		set, args = append(set, "`status` = ?"), append(args, *v)
	}
	if v := md.NextTime; v != nil {
		set, args = append(set, "`next_time` = ?"), append(args, *v)
	}
	args = append(args, md.Id)
	query := "update reminds set " + strings.Join(set, ", ") + " where id = ?"
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func (s *Store) DeleteRemind(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, "delete from reminds where id = ?", id)
	return err
}
