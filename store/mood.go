package store

import (
	"context"
	"strings"

	"app/store/model"
)

func (s *Store) ListMood(ctx context.Context, start int, num int) ([]model.Mood, error) {
	offset := (start - 1) * num
	rows, err := s.db.QueryContext(ctx, "select id,content,user_id,created_at from moods order by id desc limit ? offset ?", num, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ms := make([]model.Mood, 0)
	for rows.Next() {
		var md model.Mood
		if err := rows.Scan(&md.Id, &md.Content, &md.UserId, &md.CreatedAt); err != nil {
			return nil, err
		}
		ms = append(ms, md)
	}
	return ms, nil
}

func (s *Store) RandomMood(ctx context.Context) (*model.Mood, error) {
	var md model.Mood
	err := s.db.QueryRowContext(ctx, "select id,content,user_id,created_at from moods order by rand() limit 1").
		Scan(&md.Id, &md.Content, &md.UserId, &md.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &md, nil
}

func (s *Store) CountMoodTotal(ctx context.Context) (int, error) {
	var total int
	err := s.db.QueryRowContext(ctx, "select count(*) from moods").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (s *Store) CreateMood(ctx context.Context, md *model.Mood) (int64, error) {
	res, err := s.db.ExecContext(ctx, "insert into moods (content,user_id,created_at) values (?,?,?)", md.Content, md.UserId, md.CreatedAt)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) UpdateMood(ctx context.Context, md *model.UpdateMood) error {
	set := make([]string, 0)
	args := make([]any, 0)
	if v := md.Content; v != nil {
		set, args = append(set, "`content` = ?"), append(args, *v)
	}
	args = append(args, md.Id)
	query := "update moods set " + strings.Join(set, ", ") + " where id = ?"
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func (s *Store) DeleteMood(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, "delete from moods where id = ?", id)
	return err
}
