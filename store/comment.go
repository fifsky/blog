package store

import (
	"context"

	"app/pkg/sqlext"
	"app/store/model"
)

// commentColumns comments 表查询列
const commentColumns = "id, post_id, pid, name, email, website, reply_name, content, ip, created_at"

// commentWithPostColumns 评论关联文章标题与缩略名的查询列
const commentWithPostColumns = "c.id, c.post_id, c.pid, c.name, c.email, c.website, c.reply_name, c.content, c.ip, c.created_at," +
	" p.title as post_title, p.url as post_url"

// ListComments 查询某篇文章的全部评论，按时间正序返回（前端按 pid 分组渲染两级嵌套）
func (s *Store) ListComments(ctx context.Context, postId int) ([]model.Comment, error) {
	q := sqlext.NewBuilder().
		Select(commentColumns).
		From("comments").
		Where("post_id = ?", postId).
		OrderBy("created_at asc, id asc")

	return sqlext.Query[model.Comment](ctx, s.db, q.SQL(), q.Args()...)
}

// ListNewComments 查询最新评论（关联文章标题用于侧边栏跳转），按时间倒序
func (s *Store) ListNewComments(ctx context.Context, num int) ([]model.CommentWithPost, error) {
	q := sqlext.NewBuilder().
		Select(commentWithPostColumns).
		From("comments c left join posts p on c.post_id = p.id").
		OrderBy("c.created_at desc, c.id desc").
		Limit(num)

	return sqlext.Query[model.CommentWithPost](ctx, s.db, q.SQL(), q.Args()...)
}

// CreateComment 插入评论
func (s *Store) CreateComment(ctx context.Context, c *model.Comment) (int64, error) {
	res, err := s.db.ExecContext(ctx, "insert into comments (post_id,pid,name,email,website,reply_name,content,ip,created_at) values (?,?,?,?,?,?,?,?,?)",
		c.PostId, c.Pid, c.Name, c.Email, c.Website, c.ReplyName, c.Content, c.IP, c.CreatedAt)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// ListAllComments 后台分页查询评论（关联文章标题），支持按昵称或内容搜索
func (s *Store) ListAllComments(ctx context.Context, keyword string, start int, num int) ([]model.CommentWithPost, error) {
	q := sqlext.NewBuilder().
		Select(commentWithPostColumns).
		From("comments c left join posts p on c.post_id = p.id").
		WhereIf(keyword != "", "(c.name like ? or c.content like ?)", "%"+keyword+"%", "%"+keyword+"%").
		OrderBy("c.created_at desc, c.id desc").
		Limit(num).
		Offset((start - 1) * num)

	return sqlext.Query[model.CommentWithPost](ctx, s.db, q.SQL(), q.Args()...)
}

// CountComments 统计评论总数（支持按昵称或内容搜索）
func (s *Store) CountComments(ctx context.Context, keyword string) (int, error) {
	q := sqlext.NewBuilder().
		Select("count(*)").
		From("comments").
		WhereIf(keyword != "", "(name like ? or content like ?)", "%"+keyword+"%", "%"+keyword+"%")

	return sqlext.QueryRow[int](ctx, s.db, q.SQL(), q.Args()...)
}

// DeleteComment 根据ID批量删除评论
func (s *Store) DeleteComment(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return nil
	}
	placeholders, args := sqlext.In(ids)
	query := "delete from comments where id in (" + placeholders + ")"
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}
