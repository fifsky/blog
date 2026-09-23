package store

import (
	"context"
	"strings"
	"time"

	"app/pkg/sqlext"
	"app/store/model"
)

// remindColumns reminds 表查询列
const remindColumns = "id, cron, content, status, next_time, created_at, updated_at"

// GetRemind 按 ID 查询提醒，不存在时返回 sql.ErrNoRows
func (s *Store) GetRemind(ctx context.Context, id int) (*model.Remind, error) {
	remind, err := sqlext.QueryRow[model.Remind](ctx, s.db, "select "+remindColumns+" from reminds where id = ?", id)
	if err != nil {
		return nil, err
	}
	return &remind, nil
}

// ListRemind 分页查询提醒，按 ID 倒序
func (s *Store) ListRemind(ctx context.Context, start int, num int) ([]*model.Remind, error) {
	q := sqlext.NewBuilder().
		Select(remindColumns).
		From("reminds").
		OrderBy("id desc").
		Limit(num).
		Offset((start - 1) * num)

	list, err := sqlext.Query[model.Remind](ctx, s.db, q.SQL(), q.Args()...)
	if err != nil {
		return nil, err
	}
	return remindPointers(list), nil
}

// RemindAll 查询所有未完成的提醒（ACTIVE/PENDING），按 ID 倒序
func (s *Store) RemindAll(ctx context.Context) ([]model.Remind, error) {
	q := sqlext.NewBuilder().
		Select(remindColumns).
		From("reminds").
		Where("status in (?)", []model.RemindStatus{model.RemindStatusActive, model.RemindStatusPending}).
		OrderBy("id desc")

	return sqlext.Query[model.Remind](ctx, s.db, q.SQL(), q.Args()...)
}

// CountRemindTotal 统计提醒总数
func (s *Store) CountRemindTotal(ctx context.Context) (int, error) {
	q := sqlext.NewBuilder().Select("count(*)").From("reminds")
	return sqlext.QueryRow[int](ctx, s.db, q.SQL(), q.Args()...)
}

// UpdateRemindStatus 更新提醒状态
func (s *Store) UpdateRemindStatus(ctx context.Context, id int, status model.RemindStatus) error {
	_, err := s.db.ExecContext(ctx, "update reminds set status = ? where id = ?", status, id)
	return err
}

// UpdateRemindNextTime 更新下次提醒时间
func (s *Store) UpdateRemindNextTime(ctx context.Context, id int, nextTime time.Time) error {
	_, err := s.db.ExecContext(ctx, "update reminds set next_time = ? where id = ?", nextTime, id)
	return err
}

// CreateRemind 新增提醒，返回自增 ID
func (s *Store) CreateRemind(ctx context.Context, md *model.Remind) (int64, error) {
	res, err := s.db.ExecContext(ctx, "insert into reminds (cron,content,status,next_time,created_at,updated_at) values (?,?,?,?,?,?)",
		md.Cron, md.Content, md.Status, md.NextTime, md.CreatedAt, md.UpdatedAt)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateRemind 按需更新提醒字段，只更新传入的非空字段
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

// DeleteRemind 按 ID 删除提醒
func (s *Store) DeleteRemind(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, "delete from reminds where id = ?", id)
	return err
}

// remindPointers 把提醒值切片转成指针切片，空结果返回空切片而不是 nil
func remindPointers(list []model.Remind) []*model.Remind {
	ret := make([]*model.Remind, 0, len(list))
	for i := range list {
		ret = append(ret, &list[i])
	}
	return ret
}
