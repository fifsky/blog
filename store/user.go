package store

import (
	"context"
	"strings"

	"app/pkg/sqlext"
	"app/store/model"
)

// userColumns users 表查询列，type 是 SQL 保留字必须加反引号
const userColumns = "id, name, password, nick_name, email, status, `type`, totp_secret, created_at, updated_at"

// GetUser 按 ID 查询用户，不存在时返回 sql.ErrNoRows
func (s *Store) GetUser(ctx context.Context, uid int) (*model.User, error) {
	user, err := sqlext.QueryRow[model.User](ctx, s.db, "select "+userColumns+" from users where id = ?", uid)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ListUser 分页查询用户，按 ID 倒序
func (s *Store) ListUser(ctx context.Context, start int, num int) ([]model.User, error) {
	q := sqlext.NewBuilder().
		Select(userColumns).
		From("users").
		OrderBy("id desc").
		Limit(num).
		Offset(max((start-1)*num, 0))

	return sqlext.Query[model.User](ctx, s.db, q.SQL(), q.Args()...)
}

// CountUserTotal 统计用户总数
func (s *Store) CountUserTotal(ctx context.Context) (int, error) {
	q := sqlext.NewBuilder().Select("count(*)").From("users")
	return sqlext.QueryRow[int](ctx, s.db, q.SQL(), q.Args()...)
}

// GetUserByName 按用户名查询用户，不存在时返回 sql.ErrNoRows
func (s *Store) GetUserByName(ctx context.Context, name string) (*model.User, error) {
	q := sqlext.NewBuilder().Select(userColumns).From("users").Where("name = ?", name).Limit(1)

	user, err := sqlext.QueryRow[model.User](ctx, s.db, q.SQL(), q.Args()...)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser 新增用户，返回自增 ID
func (s *Store) CreateUser(ctx context.Context, users *model.User) (int64, error) {
	res, err := s.db.ExecContext(ctx, "insert into users (name,password,nick_name,email,status,type,totp_secret,created_at,updated_at) values (?,?,?,?,?,?,?,?,?)",
		users.Name, users.Password, users.NickName, users.Email, users.Status, users.Type, users.TotpSecret, users.CreatedAt, users.UpdatedAt)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateUser 按需更新用户字段，只更新传入的非空字段
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
	if v := users.TotpSecret; v != nil {
		set, args = append(set, "`totp_secret` = ?"), append(args, *v)
	}
	if v := users.UpdatedAt; v != nil {
		set, args = append(set, "`updated_at` = ?"), append(args, *v)
	}
	args = append(args, users.Id)
	query := "update users set " + strings.Join(set, ", ") + " where id = ?"
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

// GetUserByIds 按 ID 列表查询用户，返回以用户 ID 为键的map
func (s *Store) GetUserByIds(ctx context.Context, ids []int) (map[int]model.User, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	q := sqlext.NewBuilder().Select(userColumns).From("users").Where("id in (?)", ids)
	list, err := sqlext.Query[model.User](ctx, s.db, q.SQL(), q.Args()...)
	if err != nil {
		return nil, err
	}

	um := make(map[int]model.User, len(list))
	for _, user := range list {
		um[user.Id] = user
	}
	return um, nil
}

// GetUserIDByOpenid 按小程序 openid 查询用户 ID，不存在时返回 sql.ErrNoRows
func (s *Store) GetUserIDByOpenid(ctx context.Context, openid string) (int, error) {
	q := sqlext.NewBuilder().Select("id").From("users").Where("openid = ?", openid).Limit(1)
	return sqlext.QueryRow[int](ctx, s.db, q.SQL(), q.Args()...)
}
