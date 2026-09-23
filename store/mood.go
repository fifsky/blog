package store

import (
	"context"
	"strings"

	"app/pkg/sqlext"
	"app/store/model"
)

// moodColumns moods 表查询列
const moodColumns = "id, content, user_id, created_at, updated_at"

// ListMood 分页查询心情，按 ID 倒序
func (s *Store) ListMood(ctx context.Context, start int, num int) ([]model.Mood, error) {
	q := sqlext.NewBuilder().
		Select(moodColumns).
		From("moods").
		OrderBy("id desc").
		Limit(num).
		Offset((start - 1) * num)

	return sqlext.Query[model.Mood](ctx, s.db, q.SQL(), q.Args()...)
}

// RandomMood 随机取一条心情，空表返回 sql.ErrNoRows
func (s *Store) RandomMood(ctx context.Context) (*model.Mood, error) {
	q := sqlext.NewBuilder().Select(moodColumns).From("moods").OrderBy("random()").Limit(1)

	mood, err := sqlext.QueryRow[model.Mood](ctx, s.db, q.SQL(), q.Args()...)
	if err != nil {
		return nil, err
	}
	return &mood, nil
}

// CountMoodTotal 统计心情总数
func (s *Store) CountMoodTotal(ctx context.Context) (int, error) {
	q := sqlext.NewBuilder().Select("count(*)").From("moods")
	return sqlext.QueryRow[int](ctx, s.db, q.SQL(), q.Args()...)
}

// CreateMood 新增心情，返回自增 ID
func (s *Store) CreateMood(ctx context.Context, md *model.Mood) (int64, error) {
	res, err := s.db.ExecContext(ctx, "insert into moods (content,user_id,created_at,updated_at) values (?,?,?,?)", md.Content, md.UserId, md.CreatedAt, md.UpdatedAt)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateMood 按需更新心情内容，只更新传入的非空字段
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

// DeleteMood 根据ID批量删除心情
func (s *Store) DeleteMood(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return nil
	}
	placeholders, args := sqlext.In(ids)
	query := "delete from moods where id in (" + placeholders + ")"
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}
