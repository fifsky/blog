package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"app/store/model"
)

func (s *Store) GetRemind(ctx context.Context, id int) (*model.Remind, error) {
	query := "select id,type,content,month,week,day,hour,minute,status,next_time,created_at from blog.reminds where id = $1"
	var m model.Remind
	err := s.db.QueryRowContext(ctx, query, id).Scan(&m.Id, &m.Type, &m.Content, &m.Month, &m.Week, &m.Day, &m.Hour, &m.Minute, &m.Status, &m.NextTime, &m.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *Store) ListRemind(ctx context.Context, start int, num int) ([]*model.Remind, error) {
	offset := (start - 1) * num
	rows, err := s.db.QueryContext(ctx, "select id,type,content,month,week,day,hour,minute,status,next_time,created_at from blog.reminds order by id desc limit $1 offset $2", num, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret := make([]*model.Remind, 0)
	for rows.Next() {
		var item model.Remind
		if err := rows.Scan(&item.Id, &item.Type, &item.Content, &item.Month, &item.Week, &item.Day, &item.Hour, &item.Minute, &item.Status, &item.NextTime, &item.CreatedAt); err != nil {
			return nil, err
		}
		tmp := item
		ret = append(ret, &tmp)
	}
	return ret, nil
}

func (s *Store) RemindAll(ctx context.Context) ([]model.Remind, error) {
	rows, err := s.db.QueryContext(ctx, "select id,type,content,month,week,day,hour,minute,status,next_time,created_at from blog.reminds order by id desc")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret := make([]model.Remind, 0)
	for rows.Next() {
		var item model.Remind
		if err := rows.Scan(&item.Id, &item.Type, &item.Content, &item.Month, &item.Week, &item.Day, &item.Hour, &item.Minute, &item.Status, &item.NextTime, &item.CreatedAt); err != nil {
			return nil, err
		}
		ret = append(ret, item)
	}
	return ret, nil
}

func (s *Store) CountRemindTotal(ctx context.Context) (int, error) {
	var total int
	err := s.db.QueryRowContext(ctx, "select count(*) from blog.reminds").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (s *Store) UpdateRemindStatus(ctx context.Context, id int, status int) error {
	_, err := s.db.ExecContext(ctx, "update blog.reminds set status = $1 where id = $2", status, id)
	return err
}

func (s *Store) UpdateRemindNextTime(ctx context.Context, id int, nextTime time.Time) error {
	_, err := s.db.ExecContext(ctx, "update blog.reminds set next_time = $1 where id = $2", nextTime, id)
	return err
}

func (s *Store) CreateRemind(ctx context.Context, md *model.Remind) (int64, error) {
	var id int64
	err := s.db.QueryRowContext(ctx, "insert into blog.reminds (type,content,month,week,day,hour,minute,status,created_at) values ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id",
		md.Type, md.Content, md.Month, md.Week, md.Day, md.Hour, md.Minute, md.Status, md.CreatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Store) UpdateRemind(ctx context.Context, md *model.UpdateRemind) error {
	set := make([]string, 0)
	args := make([]any, 0)

	if v := md.Type; v != nil {
		set, args = append(set, "type = "+placeholder(len(args)+1)), append(args, *v)
	}
	if v := md.Content; v != nil {
		set, args = append(set, "content = "+placeholder(len(args)+1)), append(args, *v)
	}
	if v := md.Month; v != nil {
		set, args = append(set, "month = "+placeholder(len(args)+1)), append(args, *v)
	}
	if v := md.Week; v != nil {
		set, args = append(set, "week = "+placeholder(len(args)+1)), append(args, *v)
	}
	if v := md.Day; v != nil {
		set, args = append(set, "day = "+placeholder(len(args)+1)), append(args, *v)
	}
	if v := md.Hour; v != nil {
		set, args = append(set, "hour = "+placeholder(len(args)+1)), append(args, *v)
	}
	if v := md.Minute; v != nil {
		set, args = append(set, "minute = "+placeholder(len(args)+1)), append(args, *v)
	}
	if v := md.Status; v != nil {
		set, args = append(set, "status = "+placeholder(len(args)+1)), append(args, *v)
	}
	if v := md.NextTime; v != nil {
		set, args = append(set, "next_time = "+placeholder(len(args)+1)), append(args, *v)
	}
	args = append(args, md.Id)
	query := fmt.Sprintf("update blog.reminds set %s where id = %s", strings.Join(set, ", "), placeholder(len(args)))
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func (s *Store) DeleteRemind(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, "delete from blog.reminds where id = $1", id)
	return err
}
