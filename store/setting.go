package store

import (
	"context"
	"maps"
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
	s.optionsMu.RLock()
	if s.optionsCache != nil {
		cacheCopy := maps.Clone(s.optionsCache)
		s.optionsMu.RUnlock()
		return cacheCopy, nil
	}
	s.optionsMu.RUnlock()

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

	s.optionsMu.Lock()
	s.optionsCache = options2
	s.optionsMu.Unlock()

	// 既然已经把 options2 交给了缓存，为了防止外部修改它，
	// 返回给调用方时我们拷贝一份返回
	return maps.Clone(options2), nil
}

func (s *Store) UpdateOptions(ctx context.Context, m map[string]string) (map[string]string, error) {
	for k, v := range m {
		_, err := s.db.ExecContext(ctx, "insert into options (option_key, option_value) values (?, ?) on conflict(option_key) do update set option_value = excluded.option_value", k, v)
		if err != nil {
			return nil, err
		}
	}

	s.optionsMu.Lock()
	s.optionsCache = nil
	s.optionsMu.Unlock()

	return m, nil
}
