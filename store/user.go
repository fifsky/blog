package store

import (
	"context"
	"fmt"
	"strings"

	"app/store/model"
)

func (s *Store) GetUser(ctx context.Context, uid int) (*model.User, error) {
	query := "select id,name,password,nick_name,email,status,type,created_at,updated_at from blog.users where id = $1"
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
	query := "select id,name,password,nick_name,email,status,type,created_at,updated_at from blog.users order by id desc limit $1 offset $2"
	rows, err := s.db.QueryContext(ctx, query, num, max((start-1)*num, 0))
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
	var total int
	err := s.db.QueryRowContext(ctx, "select count(*) from blog.users").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (s *Store) GetUserByName(ctx context.Context, name string) (*model.User, error) {
	var user model.User
	err := s.db.QueryRowContext(ctx, "select id,name,password,nick_name,email,status,type,created_at,updated_at from blog.users where name = $1 limit 1", name).Scan(&user.Id, &user.Name, &user.Password, &user.NickName, &user.Email, &user.Status, &user.Type, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) CreateUser(ctx context.Context, users *model.User) (int64, error) {
	var id int64
	err := s.db.QueryRowContext(ctx, "insert into blog.users (name,password,nick_name,email,status,type,created_at,updated_at) values ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id",
		users.Name, users.Password, users.NickName, users.Email, users.Status, users.Type, users.CreatedAt, users.UpdatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Store) UpdateUser(ctx context.Context, users *model.UpdateUser) error {
	set := make([]string, 0)
	args := make([]any, 0)

	if v := users.Name; v != nil {
		set, args = append(set, "name = "+placeholder(len(args)+1)), append(args, *v)
	}
	if v := users.Password; v != nil {
		set, args = append(set, "password = "+placeholder(len(args)+1)), append(args, *v)
	}
	if v := users.NickName; v != nil {
		set, args = append(set, "nick_name = "+placeholder(len(args)+1)), append(args, *v)
	}
	if v := users.Email; v != nil {
		set, args = append(set, "email = "+placeholder(len(args)+1)), append(args, *v)
	}
	if v := users.Status; v != nil {
		set, args = append(set, "status = "+placeholder(len(args)+1)), append(args, *v)
	}
	if v := users.Type; v != nil {
		set, args = append(set, "type = "+placeholder(len(args)+1)), append(args, *v)
	}
	if v := users.UpdatedAt; v != nil {
		set, args = append(set, "updated_at = "+placeholder(len(args)+1)), append(args, *v)
	}
	args = append(args, users.Id)
	query := fmt.Sprintf("update blog.users set %s where id = %s", strings.Join(set, ", "), placeholder(len(args)))
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func (s *Store) GetUserByIds(ctx context.Context, ids []int) (map[int]model.User, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	ph, args := In(ids, 1)
	query := "select id,name,password,nick_name,email,status,type,created_at,updated_at from blog.users where id in(" + ph + ")"
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
