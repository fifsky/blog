package store

import (
	"context"
	"strings"

	"app/pkg/sqlext"
	"app/store/model"
)

// postColumns posts 表完整查询列
const postColumns = "id, cate_id, type, user_id, title, url, content, tags, status, view_num, created_at, updated_at"

// postBriefColumns 上下篇导航使用的查询列
const postBriefColumns = "id, cate_id, type, user_id, title, url, content, status, created_at, updated_at"

// GetPost 按 ID 或缩略名查询文章，ID 大于 0 时优先按 ID，不存在时返回 sql.ErrNoRows
func (s *Store) GetPost(ctx context.Context, id int, url string) (*model.Post, error) {
	q := sqlext.NewBuilder().Select(postColumns).From("posts").Limit(1)
	if id > 0 {
		q.Where("id = ?", id)
	} else {
		q.Where("url = ?", url)
	}

	p, err := sqlext.QueryRow[model.Post](ctx, s.db, q.SQL(), q.Args()...)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// IncrementPostViewNum 浏览次数加一
func (s *Store) IncrementPostViewNum(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, "update posts set view_num = view_num + 1 where id = ?", id)
	return err
}

// GetPostDaysInMonth 查询某年某月有文章的日期（去重后按日期返回）
func (s *Store) GetPostDaysInMonth(ctx context.Context, year, month int) ([]int32, error) {
	q := sqlext.NewBuilder().
		Select("distinct cast(strftime('%d', substr(created_at, 1, 19)) as integer)").
		From("posts").
		Where("status = ?", model.PostStatusActive).
		Where("cast(strftime('%Y', substr(created_at, 1, 19)) as integer) = ?", year).
		Where("cast(strftime('%m', substr(created_at, 1, 19)) as integer) = ?", month)

	return sqlext.Query[int32](ctx, s.db, q.SQL(), q.Args()...)
}

// PrevPost 查询时间上的上一篇已发布文章，不存在时返回 sql.ErrNoRows
func (s *Store) PrevPost(ctx context.Context, id int) (*model.Post, error) {
	q := sqlext.NewBuilder().
		Select(postBriefColumns).
		From("posts").
		Where("(created_at > (select created_at from posts p2 where p2.id = ?) or (created_at = (select created_at from posts p3 where p3.id = ?) and id > ?)) and status = ?",
			id, id, id, model.PostStatusActive).
		OrderBy("created_at asc, id asc").
		Limit(1)

	p, err := sqlext.QueryRow[model.Post](ctx, s.db, q.SQL(), q.Args()...)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// NextPost 查询时间上的下一篇已发布文章，不存在时返回 sql.ErrNoRows
func (s *Store) NextPost(ctx context.Context, id int) (*model.Post, error) {
	q := sqlext.NewBuilder().
		Select(postBriefColumns).
		From("posts").
		Where("(created_at < (select created_at from posts p2 where p2.id = ?) or (created_at = (select created_at from posts p3 where p3.id = ?) and id < ?)) and status = ?",
			id, id, id, model.PostStatusActive).
		OrderBy("created_at desc, id desc").
		Limit(1)

	p, err := sqlext.QueryRow[model.Post](ctx, s.db, q.SQL(), q.Args()...)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// PostArchive 按月统计已发布文章数量，按年月倒序
func (s *Store) PostArchive(ctx context.Context) ([]model.PostArchive, error) {
	q := sqlext.NewBuilder().
		Select("ym, count(ym) total").
		From("(select strftime('%Y/%m', substr(created_at, 1, 19)) as ym from posts where status = ?) s", model.PostStatusActive).
		GroupBy("ym").
		OrderBy("ym desc")

	return sqlext.Query[model.PostArchive](ctx, s.db, q.SQL(), q.Args()...)
}

// ListPost 分页查询已发布文章，支持分类/类型/日期/关键字/标签过滤
func (s *Store) ListPost(ctx context.Context, p *model.Post, start int, num int, artdate, keyword, tag string) ([]model.Post, error) {
	q := postListBuilder(p, artdate, keyword, tag).
		Select(postColumns).
		OrderBy("created_at desc, id desc").
		Limit(num).
		Offset((start - 1) * num)

	return sqlext.Query[model.Post](ctx, s.db, q.SQL(), q.Args()...)
}

// CountPosts 统计已发布文章数量，过滤条件与 ListPost 一致
func (s *Store) CountPosts(ctx context.Context, p *model.Post, artdate, keyword, tag string) (int, error) {
	q := postListBuilder(p, artdate, keyword, tag).Select("count(*)")

	return sqlext.QueryRow[int](ctx, s.db, q.SQL(), q.Args()...)
}

// postListBuilder 组装前台文章列表的过滤条件
func postListBuilder(p *model.Post, artdate, keyword, tag string) *sqlext.Builder {
	q := sqlext.NewBuilder().
		From("posts").
		Where("status = ?", model.PostStatusActive).
		WhereIf(p.CateId > 0, "cate_id = ?", p.CateId).
		WhereIf(p.Type > 0, "type = ?", p.Type)

	if artdate != "" {
		// 长度为 7 时按月过滤（YYYY-MM），否则按天过滤（YYYY-MM-DD）
		if len(artdate) == 7 {
			q.Where("strftime('%Y-%m', substr(created_at, 1, 19)) = ?", artdate)
		} else {
			q.Where("strftime('%Y-%m-%d', substr(created_at, 1, 19)) = ?", artdate)
		}
	}

	return q.
		WhereIf(keyword != "", "title like ?", "%"+keyword+"%").
		WhereIf(tag != "", "exists (select 1 from json_each(tags) where value = ?)", tag)
}

// GetCateByDomain 按域名查询分类，不存在时返回 sql.ErrNoRows
func (s *Store) GetCateByDomain(ctx context.Context, domain string) (*model.Cate, error) {
	q := sqlext.NewBuilder().Select(cateColumns).From("cates").Where("domain = ?", domain).Limit(1)

	cate, err := sqlext.QueryRow[model.Cate](ctx, s.db, q.SQL(), q.Args()...)
	if err != nil {
		return nil, err
	}
	return &cate, nil
}

// CreatePost 新增文章，返回自增 ID
func (s *Store) CreatePost(ctx context.Context, p *model.Post) (int64, error) {
	res, err := s.db.ExecContext(ctx, "insert into posts (cate_id,type,user_id,title,url,content,tags,status,created_at,updated_at) values (?,?,?,?,?,?,?,?,?,?)",
		p.CateId, p.Type, p.UserId, p.Title, p.Url, p.Content, p.Tags, p.Status, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdatePost 按需更新文章字段，只更新传入的非空字段
func (s *Store) UpdatePost(ctx context.Context, p *model.UpdatePost) error {
	set := make([]string, 0)
	args := make([]any, 0)
	if v := p.CateId; v != nil {
		set, args = append(set, "`cate_id` = ?"), append(args, *v)
	}
	if v := p.Type; v != nil {
		set, args = append(set, "`type` = ?"), append(args, *v)
	}
	if v := p.Title; v != nil {
		set, args = append(set, "`title` = ?"), append(args, *v)
	}
	if v := p.Url; v != nil {
		set, args = append(set, "`url` = ?"), append(args, *v)
	}
	if v := p.Content; v != nil {
		set, args = append(set, "`content` = ?"), append(args, *v)
	}
	if v := p.Tags; v != nil {
		set, args = append(set, "`tags` = ?"), append(args, *v)
	}
	if v := p.Status; v != nil {
		set, args = append(set, "`status` = ?"), append(args, *v)
	}
	if v := p.UpdatedAt; v != nil {
		set, args = append(set, "`updated_at` = ?"), append(args, *v)
	}
	if len(set) == 0 {
		return nil
	}
	args = append(args, p.Id)
	query := "update posts set " + strings.Join(set, ", ") + " where id = ?"
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

// SoftDeletePost 根据ID列表软删除文章（状态改为 DELETED）
func (s *Store) SoftDeletePost(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return nil
	}
	placeholders, args := sqlext.In(ids)
	query := "update posts set status = 'DELETED' where id in (" + placeholders + ")"
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

// RestorePost 把软删除的文章恢复为草稿
func (s *Store) RestorePost(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, "update posts set status = 'DRAFT' where id = ?", id)
	return err
}

// ListPostForAdmin 后台分页查询文章，支持分类/类型/状态/关键字过滤
func (s *Store) ListPostForAdmin(ctx context.Context, p *model.Post, start int, num int, keyword string) ([]model.Post, error) {
	q := postAdminBuilder(p, keyword).
		Select(postColumns).
		OrderBy("created_at desc, id desc").
		Limit(num).
		Offset((start - 1) * num)

	return sqlext.Query[model.Post](ctx, s.db, q.SQL(), q.Args()...)
}

// CountPostsForAdmin 统计后台文章数量，过滤条件与 ListPostForAdmin 一致
func (s *Store) CountPostsForAdmin(ctx context.Context, p *model.Post, keyword string) (int, error) {
	q := postAdminBuilder(p, keyword).Select("count(*)")

	return sqlext.QueryRow[int](ctx, s.db, q.SQL(), q.Args()...)
}

// postAdminBuilder 组装后台文章列表的过滤条件
func postAdminBuilder(p *model.Post, keyword string) *sqlext.Builder {
	return sqlext.NewBuilder().
		From("posts").
		WhereIf(p.CateId > 0, "cate_id = ?", p.CateId).
		WhereIf(p.Type > 0, "type = ?", p.Type).
		WhereIf(p.Status != "", "status = ?", p.Status).
		WhereIf(keyword != "", "title like ?", "%"+keyword+"%")
}

// DestroyPost 根据ID列表物理删除已软删除的文章
func (s *Store) DestroyPost(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return nil
	}
	placeholders, args := sqlext.In(ids)
	query := "delete from posts where status = 'DELETED' and id in (" + placeholders + ")"
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}
