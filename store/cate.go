package store

import (
	"context"
	"fmt"
	"strings"

	"app/store/model"
)

func (s *Store) GetCate(ctx context.Context, id int) (*model.Cate, error) {
	row := s.db.QueryRowContext(ctx, `select id,name,"desc",domain,created_at,updated_at from blog.cates where id = $1`, id)
	c := model.Cate{}
	if err := row.Scan(&c.Id, &c.Name, &c.Desc, &c.Domain, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *Store) GetAllCates(ctx context.Context) ([]model.CateArtivleCount, error) {
	rows, err := s.db.QueryContext(ctx, "select c.id,c.name,c.desc,c.domain,c.created_at,c.updated_at,COALESCE(p.num,0) num from blog.cates c left join (select count(*) num ,cate_id from blog.posts where status = 1 and type = 1 group by cate_id) p on c.id = p.cate_id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cs []model.CateArtivleCount
	for rows.Next() {
		c := model.CateArtivleCount{}
		if err := rows.Scan(&c.Id, &c.Name, &c.Desc, &c.Domain, &c.CreatedAt, &c.UpdatedAt, &c.Num); err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, nil
}

func (s *Store) CreateCate(ctx context.Context, c *model.Cate) (int64, error) {
	var id int64
	err := s.db.QueryRowContext(ctx, `insert into blog.cates (name,"desc",domain,created_at,updated_at) values ($1,$2,$3,$4,$5) RETURNING id`,
		c.Name, c.Desc, c.Domain, c.CreatedAt, c.UpdatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Store) UpdateCate(ctx context.Context, c *model.UpdateCate) error {
	var (
		set  []string
		args []any
	)

	if v := c.Name; v != nil {
		set, args = append(set, "name = "+placeholder(len(args)+1)), append(args, *v)
	}
	if v := c.Desc; v != nil {
		set, args = append(set, `"desc" = `+placeholder(len(args)+1)), append(args, *v)
	}
	if v := c.Domain; v != nil {
		set, args = append(set, "domain = "+placeholder(len(args)+1)), append(args, *v)
	}
	if v := c.UpdatedAt; v != nil {
		set, args = append(set, "updated_at = "+placeholder(len(args)+1)), append(args, *v)
	}
	args = append(args, c.Id)

	query := fmt.Sprintf("UPDATE cates SET %s WHERE id = %s", strings.Join(set, ", "), placeholder(len(args)))
	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	return nil
}

func (s *Store) DeleteCate(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, "delete from blog.cates where id = $1", id)
	return err
}

func (s *Store) GetCatesByIds(ctx context.Context, ids []int) (map[int]model.Cate, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	ph, args := In(ids, 1)
	query := `select id,name,"desc",domain,created_at,updated_at from blog.cates where id in(` + ph + `)`
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cm := make(map[int]model.Cate)
	for rows.Next() {
		var c model.Cate
		if err := rows.Scan(&c.Id, &c.Name, &c.Desc, &c.Domain, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		cm[c.Id] = c
	}
	return cm, nil
}

func (s *Store) PostsCount(ctx context.Context, cateId int) (int, error) {
	var total int
	err := s.db.QueryRowContext(ctx, "select count(*) from blog.posts where cate_id = $1", cateId).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}
