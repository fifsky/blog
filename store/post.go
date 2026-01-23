package store

import (
	"context"
	"fmt"
	"strings"

	"app/store/model"
)

func (s *Store) GetPost(ctx context.Context, id int, url string) (*model.Post, error) {
	var p model.Post

	var where string
	var arg any
	if id > 0 {
		where = "id = ?"
		arg = id
	} else {
		where = "url = ?"
		arg = url
	}

	query := "select id,cate_id,type,user_id,title,url,content,tags,status,view_num,created_at,updated_at from posts where " + where + " limit 1"
	err := s.db.QueryRowContext(ctx, query, arg).Scan(&p.Id, &p.CateId, &p.Type, &p.UserId, &p.Title, &p.Url, &p.Content, &p.Tags, &p.Status, &p.ViewNum, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *Store) IncrementPostViewNum(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, "update posts set view_num = view_num + 1 where id = ?", id)
	return err
}

func (s *Store) GetPostDaysInMonth(ctx context.Context, year, month int) ([]int32, error) {
	query := "select distinct DAY(created_at) from posts where type = 1 and status = 1 and YEAR(created_at) = ? and MONTH(created_at) = ?"
	rows, err := s.db.QueryContext(ctx, query, year, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var days []int32
	for rows.Next() {
		var day int32
		if err := rows.Scan(&day); err != nil {
			return nil, err
		}
		days = append(days, day)
	}
	return days, nil
}

func (s *Store) PrevPost(ctx context.Context, id int) (*model.Post, error) {
	var p model.Post
	err := s.db.QueryRowContext(ctx, "select id,cate_id,type,user_id,title,url,content,status,created_at,updated_at from posts where id < ? and status = 1 and type = 1 order by id desc limit 1", id).Scan(&p.Id, &p.CateId, &p.Type, &p.UserId, &p.Title, &p.Url, &p.Content, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Store) NextPost(ctx context.Context, id int) (*model.Post, error) {
	var p model.Post
	err := s.db.QueryRowContext(ctx, "select id,cate_id,type,user_id,title,url,content,status,created_at,updated_at from posts where id > ? and status = 1 and type = 1 order by id asc limit 1", id).Scan(&p.Id, &p.CateId, &p.Type, &p.UserId, &p.Title, &p.Url, &p.Content, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Store) PostArchive(ctx context.Context) ([]model.PostArchive, error) {
	res := make([]model.PostArchive, 0)
	rows, err := s.db.QueryContext(ctx, "select ym,count(ym) total from (select DATE_FORMAT(created_at,'%Y/%m') as ym from posts where type = 1 and status = 1) s group by ym order by ym desc")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var ym, total string
		if err := rows.Scan(&ym, &total); err != nil {
			return nil, err
		}
		res = append(res, model.PostArchive{
			Ym:    ym,
			Total: total,
		})
	}
	return res, nil
}

func (s *Store) ListPost(ctx context.Context, p *model.Post, start int, num int, artdate, keyword, tag string) ([]model.Post, error) {
	posts := make([]model.Post, 0)
	offset := (start - 1) * num

	args := make([]any, 0)
	where := "status = 1"

	if p.CateId > 0 {
		where += " and cate_id = ?"
		args = append(args, p.CateId)
	}
	if p.Type > 0 {
		where += " and type = ?"
		args = append(args, p.Type)
	}
	if artdate != "" {
		if len(artdate) == 7 {
			where += " and DATE_FORMAT(created_at,'%Y-%m') = ?"
		} else {
			where += " and DATE_FORMAT(created_at,'%Y-%m-%d') = ?"
		}
		args = append(args, artdate)
	}
	if keyword != "" {
		where += " and title like ?"
		args = append(args, fmt.Sprintf("%%%s%%", keyword))
	}
	if tag != "" {
		where += " and JSON_CONTAINS(tags, JSON_QUOTE(?))"
		args = append(args, tag)
	}
	args = append(args, num, offset)

	query := "select id,cate_id,type,user_id,title,url,content,tags,status,created_at,updated_at from posts where " + where + " order by id desc limit ? offset ?"
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var bp model.Post
		if err := rows.Scan(&bp.Id, &bp.CateId, &bp.Type, &bp.UserId, &bp.Title, &bp.Url, &bp.Content, &bp.Tags, &bp.Status, &bp.CreatedAt, &bp.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, bp)
	}
	return posts, nil
}

func (s *Store) CountPosts(ctx context.Context, p *model.Post, artdate, keyword, tag string) (int, error) {
	args := make([]any, 0)
	where := "status = 1"
	if p.CateId > 0 {
		where += " and cate_id = ?"
		args = append(args, p.CateId)
	}
	if p.Type > 0 {
		where += " and type = ?"
		args = append(args, p.Type)
	}
	if artdate != "" {
		if len(artdate) == 7 {
			where += " and DATE_FORMAT(created_at,'%Y-%m') = ?"
		} else {
			where += " and DATE_FORMAT(created_at,'%Y-%m-%d') = ?"
		}
		args = append(args, artdate)
	}
	if keyword != "" {
		where += " and title like ?"
		args = append(args, fmt.Sprintf("%%%s%%", keyword))
	}
	if tag != "" {
		where += " and JSON_CONTAINS(tags, JSON_QUOTE(?))"
		args = append(args, tag)
	}
	q := "select count(*) from posts where " + where
	var total int
	err := s.db.QueryRowContext(ctx, q, args...).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (s *Store) GetCateByDomain(ctx context.Context, domain string) (*model.Cate, error) {
	var c model.Cate
	err := s.db.QueryRowContext(ctx, "select id,name,`desc`,domain,created_at,updated_at from cates where domain = ? limit 1", domain).Scan(&c.Id, &c.Name, &c.Desc, &c.Domain, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *Store) CreatePost(ctx context.Context, p *model.Post) (int64, error) {
	res, err := s.db.ExecContext(ctx, "insert into posts (cate_id,type,user_id,title,url,content,tags,status,created_at,updated_at) values (?,?,?,?,?,?,?,?,?,?)",
		p.CateId, p.Type, p.UserId, p.Title, p.Url, p.Content, p.Tags, p.Status, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

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

func (s *Store) SoftDeletePost(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return nil
	}
	placeholders, args := In(ids)
	query := "update posts set status = 2 where id in (" + placeholders + ")"
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func (s *Store) RestorePost(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, "update posts set status = 3 where id = ?", id)
	return err
}

func (s *Store) ListPostForAdmin(ctx context.Context, p *model.Post, start int, num int) ([]model.Post, error) {
	posts := make([]model.Post, 0)
	offset := (start - 1) * num

	args := make([]any, 0)
	where := "1=1"

	if p.CateId > 0 {
		where += " and cate_id = ?"
		args = append(args, p.CateId)
	}
	if p.Type > 0 {
		where += " and type = ?"
		args = append(args, p.Type)
	}
	if p.Status > 0 {
		where += " and status = ?"
		args = append(args, p.Status)
	}
	args = append(args, num, offset)

	query := "select id,cate_id,type,user_id,title,url,content,tags,status,view_num,created_at,updated_at from posts where " + where + " order by id desc limit ? offset ?"
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var bp model.Post
		if err := rows.Scan(&bp.Id, &bp.CateId, &bp.Type, &bp.UserId, &bp.Title, &bp.Url, &bp.Content, &bp.Tags, &bp.Status, &bp.ViewNum, &bp.CreatedAt, &bp.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, bp)
	}
	return posts, nil
}

func (s *Store) CountPostsForAdmin(ctx context.Context, p *model.Post) (int, error) {
	args := make([]any, 0)
	where := "1=1"
	if p.CateId > 0 {
		where += " and cate_id = ?"
		args = append(args, p.CateId)
	}
	if p.Type > 0 {
		where += " and type = ?"
		args = append(args, p.Type)
	}
	if p.Status > 0 {
		where += " and status = ?"
		args = append(args, p.Status)
	}
	q := "select count(*) from posts where " + where
	var total int
	err := s.db.QueryRowContext(ctx, q, args...).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}
