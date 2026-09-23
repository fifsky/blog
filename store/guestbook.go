package store

import (
	"context"

	"app/pkg/sqlext"
	"app/store/model"
)

// guestbookColumns guestbook 表查询列
const guestbookColumns = "id, name, content, ip, top, created_at"

// ListGuestbook 分页查询留言，支持按昵称或内容搜索，置顶优先
func (s *Store) ListGuestbook(ctx context.Context, keyword string, start int, num int) ([]model.Guestbook, error) {
	q := sqlext.NewBuilder().
		Select(guestbookColumns).
		From("guestbook").
		WhereIf(keyword != "", "(name like ? or content like ?)", "%"+keyword+"%", "%"+keyword+"%").
		OrderBy("top desc, created_at desc").
		Limit(num).
		Offset((start - 1) * num)

	return sqlext.Query[model.Guestbook](ctx, s.db, q.SQL(), q.Args()...)
}

// CountGuestbookTotal 统计留言总数（支持按昵称或内容搜索）
func (s *Store) CountGuestbookTotal(ctx context.Context, keyword string) (int, error) {
	q := sqlext.NewBuilder().
		Select("count(*)").
		From("guestbook").
		WhereIf(keyword != "", "(name like ? or content like ?)", "%"+keyword+"%", "%"+keyword+"%")

	return sqlext.QueryRow[int](ctx, s.db, q.SQL(), q.Args()...)
}

// CreateGuestbook 新增留言
func (s *Store) CreateGuestbook(ctx context.Context, g *model.Guestbook) (int64, error) {
	res, err := s.db.ExecContext(ctx, "insert into guestbook (name,content,ip,created_at) values (?,?,?,?)", g.Name, g.Content, g.Ip, g.CreatedAt)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetGuestbook 按 ID 查询留言，不存在时返回 sql.ErrNoRows
func (s *Store) GetGuestbook(ctx context.Context, id int64) (*model.Guestbook, error) {
	g, err := sqlext.QueryRow[model.Guestbook](ctx, s.db, "select "+guestbookColumns+" from guestbook where id = ?", id)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

// DeleteGuestbook 根据ID批量删除留言
func (s *Store) DeleteGuestbook(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return nil
	}
	placeholders, args := sqlext.In(ids)
	query := "delete from guestbook where id in (" + placeholders + ")"
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}
