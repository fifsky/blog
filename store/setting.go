package store

import (
	"context"
	"maps"

	"app/pkg/sqlext"
	"app/store/model"
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

// GetOptions 获取全部配置项，结果会缓存；返回给调用方的是副本，可自由修改
func (s *Store) GetOptions(ctx context.Context) (map[string]string, error) {
	s.optionsMu.RLock()
	if s.optionsCache != nil {
		cacheCopy := maps.Clone(s.optionsCache)
		s.optionsMu.RUnlock()
		return cacheCopy, nil
	}
	s.optionsMu.RUnlock()

	q := sqlext.NewBuilder().Select("option_key, option_value").From("options")

	list, err := sqlext.Query[model.Option](ctx, s.db, q.SQL(), q.Args()...)
	if err != nil {
		return nil, err
	}

	options := make(map[string]string, len(list))
	for _, opt := range list {
		options[opt.OptionKey] = opt.OptionValue
	}

	s.optionsMu.Lock()
	s.optionsCache = options
	s.optionsMu.Unlock()

	// 既然已经把 options 交给了缓存，为了防止外部修改它，
	// 返回给调用方时我们拷贝一份返回
	return maps.Clone(options), nil
}

// UpdateOptions 逐条写入配置项（存在则更新），并让缓存失效
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
