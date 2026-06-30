package store

import (
	"context"
	"strings"

	"app/store/model"
)

// ListComments 查询某篇文章的全部评论，按时间正序返回（前端按 pid 分组渲染两级嵌套）
func (s *Store) ListComments(ctx context.Context, postId int) ([]model.Comment, error) {
	query := "select id,post_id,pid,name,email,website,reply_name,content,ip,created_at from comments where post_id = ? order by created_at asc, id asc"
	rows, err := s.db.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]model.Comment, 0)
	for rows.Next() {
		var c model.Comment
		if err := rows.Scan(&c.Id, &c.PostId, &c.Pid, &c.Name, &c.Email, &c.Website, &c.ReplyName, &c.Content, &c.IP, &c.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}

// ListNewComments 查询最新评论（关联文章标题用于侧边栏跳转），按时间倒序
func (s *Store) ListNewComments(ctx context.Context, num int) ([]model.CommentWithPost, error) {
	query := `select c.id, c.post_id, c.pid, c.name, c.email, c.website, c.reply_name, c.content, c.ip, c.created_at, p.title, p.url
		from comments c left join posts p on c.post_id = p.id
		order by c.created_at desc, c.id desc limit ?`
	rows, err := s.db.QueryContext(ctx, query, num)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]model.CommentWithPost, 0)
	for rows.Next() {
		var c model.CommentWithPost
		if err := rows.Scan(&c.Id, &c.PostId, &c.Pid, &c.Name, &c.Email, &c.Website, &c.ReplyName, &c.Content, &c.IP, &c.CreatedAt, &c.PostTitle, &c.PostUrl); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
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
	offset := (start - 1) * num
	query := `select c.id, c.post_id, c.pid, c.name, c.email, c.website, c.reply_name, c.content, c.ip, c.created_at, p.title, p.url
		from comments c left join posts p on c.post_id = p.id`
	args := []any{}

	if keyword != "" {
		query += " where (c.name like ? or c.content like ?)"
		like := "%" + keyword + "%"
		args = append(args, like, like)
	}

	query += " order by c.created_at desc, c.id desc limit ? offset ?"
	args = append(args, num, offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]model.CommentWithPost, 0)
	for rows.Next() {
		var c model.CommentWithPost
		if err := rows.Scan(&c.Id, &c.PostId, &c.Pid, &c.Name, &c.Email, &c.Website, &c.ReplyName, &c.Content, &c.IP, &c.CreatedAt, &c.PostTitle, &c.PostUrl); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}

// CountComments 统计评论总数（支持按昵称或内容搜索）
func (s *Store) CountComments(ctx context.Context, keyword string) (int, error) {
	query := "select count(*) from comments"
	args := []any{}

	if keyword != "" {
		query += " where (name like ? or content like ?)"
		like := "%" + keyword + "%"
		args = append(args, like, like)
	}

	var total int
	err := s.db.QueryRowContext(ctx, query, args...).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// DeleteComment 根据ID批量删除评论
func (s *Store) DeleteComment(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return nil
	}
	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	query := "delete from comments where id in (" + strings.Join(placeholders, ",") + ")"
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}
