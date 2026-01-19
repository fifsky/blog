package store

import (
	"context"
	"strings"
	"time"

	"app/store/model"
)

func (s *Store) GetRemind(ctx context.Context, id int) (*model.Remind, error) {
	query := "select id,`type`,content,month,week,`day`,`hour`,minute,status,next_time,created_at from reminds where id = ?"
	var m model.Remind
	err := s.db.QueryRowContext(ctx, query, id).Scan(&m.Id, &m.Type, &m.Content, &m.Month, &m.Week, &m.Day, &m.Hour, &m.Minute, &m.Status, &m.NextTime, &m.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *Store) ListRemind(ctx context.Context, start int, num int) ([]*model.Remind, error) {
	offset := (start - 1) * num
	rows, err := s.db.QueryContext(ctx, "select id,`type`,content,month,week,`day`,`hour`,minute,status,next_time,created_at from reminds order by id desc limit ? offset ?", num, offset)
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
	rows, err := s.db.QueryContext(ctx, "select id,`type`,content,month,week,`day`,`hour`,minute,status,next_time,created_at from reminds order by id desc")
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
	err := s.db.QueryRowContext(ctx, "select count(*) from reminds").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (s *Store) UpdateRemindStatus(ctx context.Context, id int, status int) error {
	_, err := s.db.ExecContext(ctx, "update reminds set status = ? where id = ?", status, id)
	return err
}

func (s *Store) UpdateRemindNextTime(ctx context.Context, id int, nextTime time.Time) error {
	_, err := s.db.ExecContext(ctx, "update reminds set next_time = ? where id = ?", nextTime, id)
	return err
}

func (s *Store) CreateRemind(ctx context.Context, md *model.Remind) (int64, error) {
	res, err := s.db.ExecContext(ctx, "insert into reminds (type,content,month,week,day,hour,minute,status,created_at) values (?,?,?,?,?,?,?,?,?)",
		md.Type, md.Content, md.Month, md.Week, md.Day, md.Hour, md.Minute, md.Status, md.CreatedAt)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) UpdateRemind(ctx context.Context, md *model.UpdateRemind) error {
	set := make([]string, 0)
	args := make([]any, 0)
	if v := md.Type; v != nil {
		set, args = append(set, "`type` = ?"), append(args, *v)
	}
	if v := md.Content; v != nil {
		set, args = append(set, "`content` = ?"), append(args, *v)
	}
	if v := md.Month; v != nil {
		set, args = append(set, "`month` = ?"), append(args, *v)
	}
	if v := md.Week; v != nil {
		set, args = append(set, "`week` = ?"), append(args, *v)
	}
	if v := md.Day; v != nil {
		set, args = append(set, "`day` = ?"), append(args, *v)
	}
	if v := md.Hour; v != nil {
		set, args = append(set, "`hour` = ?"), append(args, *v)
	}
	if v := md.Minute; v != nil {
		set, args = append(set, "`minute` = ?"), append(args, *v)
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
