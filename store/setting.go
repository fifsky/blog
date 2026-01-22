package store

import (
	"context"
)

func (s *Store) GetOptions(ctx context.Context) (map[string]string, error) {
	rows, err := s.db.QueryContext(ctx, "select id,option_key,option_value from blog.options")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	options2 := make(map[string]string)
	for rows.Next() {
		var id int
		var k, v string
		if err := rows.Scan(&id, &k, &v); err != nil {
			return nil, err
		}
		options2[k] = v
	}
	return options2, nil
}

func (s *Store) UpdateOptions(ctx context.Context, m map[string]string) (map[string]string, error) {
	for k, v := range m {
		_, err := s.db.ExecContext(ctx, "update blog.options set option_value = $1 where option_key = $2", v, k)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}
