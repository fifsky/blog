package store

import (
	"context"
)

// AIConfig AI 服务配置
type AIConfig struct {
	Token    string
	Endpoint string
	Model    string
}

// GetAIConfig 获取 AI 配置，从数据库读取
func (s *Store) GetAIConfig(ctx context.Context) *AIConfig {
	opts, err := s.GetOptions(ctx)
	if err != nil {
		return &AIConfig{}
	}

	return &AIConfig{
		Token:    opts["ai_token"],
		Endpoint: opts["ai_endpoint"],
		Model:    opts["ai_model"],
	}
}

func (s *Store) GetOptions(ctx context.Context) (map[string]string, error) {
	rows, err := s.db.QueryContext(ctx, "select id,option_key,option_value from options")
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
		_, err := s.db.ExecContext(ctx, "update options set option_value = ? where option_key = ?", v, k)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}
