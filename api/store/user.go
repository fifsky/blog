package store

import (
	"context"
	"database/sql"
	"strings"

	"app/model"
)

func (s *Store) GetUser(ctx context.Context, uid int) (*model.User, error) {
	query := "select id,name,password,nick_name,email,status,`type`,created_at,updated_at from users where id = ?"
	row := s.db.QueryRowContext(ctx, query, uid)
	if row.Err() != nil {
		return nil, row.Err()
	}
	var user model.User
	if err := row.Scan(&user.Id, &user.Name, &user.Password, &user.NickName, &user.Email, &user.Status, &user.Type, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) ListUser(ctx context.Context, start int, num int) ([]model.User, error) {
	query := "select id,name,password,nick_name,email,status,`type`,created_at,updated_at from users order by id desc limit ?,?"
	rows, err := s.db.QueryContext(ctx, query, max((start-1)*num, 0), num)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]model.User, 0)
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.Id, &user.Name, &user.Password, &user.NickName, &user.Email, &user.Status, &user.Type, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *Store) CountUserTotal(ctx context.Context) (int, error) {
	rows, err := s.db.QueryContext(ctx, "select count(*) from users")
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var total int
	if rows.Next() {
		if err := rows.Scan(&total); err != nil {
			return 0, err
		}
	}
	return total, nil
}

func (s *Store) GetUserByName(ctx context.Context, name string) (*model.User, error) {
	rows, err := s.db.QueryContext(ctx, "select id,name,password,nick_name,email,status,`type`,created_at,updated_at from users where name = ? limit 1", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var user model.User
		if err := rows.Scan(&user.Id, &user.Name, &user.Password, &user.NickName, &user.Email, &user.Status, &user.Type, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		return &user, nil
	}
	return nil, sql.ErrNoRows
}

func (s *Store) CreateUser(ctx context.Context, users *model.User) (int64, error) {
	res, err := s.db.ExecContext(ctx, "insert into users (name,password,nick_name,email,status,type,created_at,updated_at) values (?,?,?,?,?,?,?,?)",
		users.Name, users.Password, users.NickName, users.Email, users.Status, users.Type, users.CreatedAt, users.UpdatedAt)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) UpdateUser(ctx context.Context, users *model.UpdateUser) error {
	set := make([]string, 0)
	args := make([]any, 0)
	if v := users.Name; v != nil {
		set, args = append(set, "`name` = ?"), append(args, *v)
	}
	if v := users.Password; v != nil {
		set, args = append(set, "`password` = ?"), append(args, *v)
	}
	if v := users.NickName; v != nil {
		set, args = append(set, "`nick_name` = ?"), append(args, *v)
	}
	if v := users.Email; v != nil {
		set, args = append(set, "`email` = ?"), append(args, *v)
	}
	if v := users.Status; v != nil {
		set, args = append(set, "`status` = ?"), append(args, *v)
	}
	if v := users.Type; v != nil {
		set, args = append(set, "`type` = ?"), append(args, *v)
	}
	if v := users.UpdatedAt; v != nil {
		set, args = append(set, "`updated_at` = ?"), append(args, *v)
	}
	args = append(args, users.Id)
	query := "update users set " + strings.Join(set, ", ") + " where id = ?"
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func (s *Store) GetUserByIds(ctx context.Context, ids []int) (map[int]model.User, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	ph, args := In(ids)
	query := "select id,name,password,nick_name,email,status,`type`,created_at,updated_at from users where id in(" + ph + ")"
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	um := make(map[int]model.User)
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.Id, &user.Name, &user.Password, &user.NickName, &user.Email, &user.Status, &user.Type, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		um[user.Id] = user
	}
	return um, nil
}
