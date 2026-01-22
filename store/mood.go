package store

import (
	"context"
	"fmt"
	"strings"

	"app/store/model"
)

func (s *Store) ListMood(ctx context.Context, start int, num int) ([]model.Mood, error) {
	offset := (start - 1) * num
	rows, err := s.db.QueryContext(ctx, "select id,content,user_id,created_at from blog.moods order by id desc limit $1 offset $2", num, offset)
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
	err := s.db.QueryRowContext(ctx, "select id,content,user_id,created_at from blog.moods order by random() limit 1").
		Scan(&md.Id, &md.Content, &md.UserId, &md.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &md, nil
}

func (s *Store) CountMoodTotal(ctx context.Context) (int, error) {
	var total int
	err := s.db.QueryRowContext(ctx, "select count(*) from blog.moods").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (s *Store) CreateMood(ctx context.Context, md *model.Mood) (int64, error) {
	var id int64
	err := s.db.QueryRowContext(ctx, "insert into blog.moods (content,user_id,created_at) values ($1,$2,$3) RETURNING id", md.Content, md.UserId, md.CreatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Store) UpdateMood(ctx context.Context, md *model.UpdateMood) error {
	set := make([]string, 0)
	args := make([]any, 0)

	if v := md.Content; v != nil {
		set, args = append(set, "content = "+placeholder(len(args)+1)), append(args, *v)
	}
	args = append(args, md.Id)
	query := fmt.Sprintf("update blog.moods set %s where id = %s", strings.Join(set, ", "), placeholder(len(args)))
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func (s *Store) DeleteMood(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, "delete from blog.moods where id = $1", id)
	return err
}
