package store

import (
	"context"
	"strings"

	"app/pkg/sqlext"
	"app/store/model"
)

// linkColumns links 表查询列，desc 是 SQL 保留字必须加反引号
const linkColumns = "id, name, url, `desc`, status, created_at, updated_at"

// GetAllLinks 查询全部友链
func (s *Store) GetAllLinks(ctx context.Context) ([]model.Link, error) {
	q := sqlext.NewBuilder().Select(linkColumns).From("links").OrderBy("id asc")

	return sqlext.Query[model.Link](ctx, s.db, q.SQL(), q.Args()...)
}

// GetApprovedLinks 获取审核通过的链接列表
func (s *Store) GetApprovedLinks(ctx context.Context) ([]model.Link, error) {
	q := sqlext.NewBuilder().
		Select(linkColumns).
		From("links").
		Where("status = ?", model.LinkStatusApproved).
		OrderBy("id asc")

	return sqlext.Query[model.Link](ctx, s.db, q.SQL(), q.Args()...)
}

// GetLink 根据 ID 获取链接信息，不存在时返回 sql.ErrNoRows
func (s *Store) GetLink(ctx context.Context, id int) (*model.Link, error) {
	link, err := sqlext.QueryRow[model.Link](ctx, s.db, "select "+linkColumns+" from links where id = ?", id)
	if err != nil {
		return nil, err
	}
	return &link, nil
}

// CreateLink 新增友链，返回自增 ID
func (s *Store) CreateLink(ctx context.Context, link *model.Link) (int64, error) {
	res, err := s.db.ExecContext(ctx, "insert into links (name,url,`desc`,status,created_at,updated_at) values (?,?,?,?,?,?)",
		link.Name, link.Url, link.Desc, link.Status, link.CreatedAt, link.UpdatedAt)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateLink 按需更新友链字段，只更新传入的非空字段
func (s *Store) UpdateLink(ctx context.Context, link *model.UpdateLink) error {
	set := make([]string, 0)
	args := make([]any, 0)
	if v := link.Name; v != nil {
		set, args = append(set, "`name` = ?"), append(args, *v)
	}
	if v := link.Url; v != nil {
		set, args = append(set, "`url` = ?"), append(args, *v)
	}
	if v := link.Desc; v != nil {
		set, args = append(set, "`desc` = ?"), append(args, *v)
	}
	if v := link.Status; v != nil {
		set, args = append(set, "`status` = ?"), append(args, *v)
	}
	args = append(args, link.Id)
	query := "update links set " + strings.Join(set, ", ") + " where id = ?"
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

// DeleteLink 按 ID 删除友链
func (s *Store) DeleteLink(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, "delete from links where id = ?", id)
	return err
}
