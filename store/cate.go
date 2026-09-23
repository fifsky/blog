package store

import (
	"context"
	"strings"

	"app/pkg/sqlext"
	"app/store/model"
)

// cateColumns cates 表查询列，desc 是 SQL 保留字必须加反引号
const cateColumns = "id, name, `desc`, domain, created_at, updated_at"

// GetCate 按 ID 查询分类，不存在时返回 sql.ErrNoRows
func (s *Store) GetCate(ctx context.Context, id int) (*model.Cate, error) {
	cate, err := sqlext.QueryRow[model.Cate](ctx, s.db, "select "+cateColumns+" from cates where id = ?", id)
	if err != nil {
		return nil, err
	}
	return &cate, nil
}

// GetAllCates 查询全部分类及每个分类下的文章数量
func (s *Store) GetAllCates(ctx context.Context) ([]model.CateArtivleCount, error) {
	query := "select c.id, c.name, c.`desc`, c.domain, c.created_at, c.updated_at, ifnull(p.num, 0) as num" +
		" from cates c left join (select count(*) num, cate_id from posts where status = 'ACTIVE' and type = 1 group by cate_id) p on c.id = p.cate_id"
	return sqlext.Query[model.CateArtivleCount](ctx, s.db, query)
}

// CreateCate 新增分类，返回自增 ID
func (s *Store) CreateCate(ctx context.Context, c *model.Cate) (int64, error) {
	res, err := s.db.ExecContext(ctx, "insert into cates (name,`desc`,domain,created_at,updated_at) values (?,?,?,?,?)",
		c.Name, c.Desc, c.Domain, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateCate 按需更新分类字段，只更新传入的非空字段
func (s *Store) UpdateCate(ctx context.Context, c *model.UpdateCate) error {
	var (
		set  []string
		args []any
	)
	if v := c.Name; v != nil {
		set, args = append(set, "`name` = ?"), append(args, *v)
	}
	if v := c.Desc; v != nil {
		set, args = append(set, "`desc` = ?"), append(args, *v)
	}
	if v := c.Domain; v != nil {
		set, args = append(set, "`domain` = ?"), append(args, *v)
	}
	if v := c.UpdatedAt; v != nil {
		set, args = append(set, "`updated_at` = ?"), append(args, *v)
	}
	args = append(args, c.Id)

	query := "UPDATE `cates` SET " + strings.Join(set, ", ") + " WHERE `id` = ?"
	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	return nil
}

// DeleteCate 按 ID 删除分类
func (s *Store) DeleteCate(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, "delete from cates where id = ?", id)
	return err
}

// GetCatesByIds 按 ID 列表查询分类，返回以分类 ID 为键的map
func (s *Store) GetCatesByIds(ctx context.Context, ids []int) (map[int]model.Cate, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	q := sqlext.NewBuilder().Select(cateColumns).From("cates").Where("id in (?)", ids)
	list, err := sqlext.Query[model.Cate](ctx, s.db, q.SQL(), q.Args()...)
	if err != nil {
		return nil, err
	}

	cm := make(map[int]model.Cate, len(list))
	for _, cate := range list {
		cm[cate.Id] = cate
	}
	return cm, nil
}

// PostsCount 统计分类下的文章数量
func (s *Store) PostsCount(ctx context.Context, cateId int) (int, error) {
	q := sqlext.NewBuilder().Select("count(*)").From("posts").Where("cate_id = ?", cateId)
	return sqlext.QueryRow[int](ctx, s.db, q.SQL(), q.Args()...)
}
