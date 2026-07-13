package store

import (
	"context"
	"testing"

	"app/pkg/dbunit"
	"app/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetting_GetOptions(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options"))
		s := New(db)
		ret, err := s.GetOptions(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		assert.Equal(t, ret["site_name"], "無處告別")
	})
}

func TestSetting_UpdateOptions(t *testing.T) {
	tests := []struct {
		name string
		kv   map[string]string
		want map[string]string
	}{
		{
			name: "更新已存在配置",
			kv:   map[string]string{"site_name": "新站点"},
			want: map[string]string{"site_name": "新站点"},
		},
		{
			name: "保存不存在的AI配置",
			kv: map[string]string{
				"ai_token":    "test-token",
				"ai_endpoint": "https://example.com/v1",
				"ai_model":    "test-model",
			},
			want: map[string]string{
				"ai_token":    "test-token",
				"ai_endpoint": "https://example.com/v1",
				"ai_model":    "test-model",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options"))
				s := New(db)

				_, err := s.UpdateOptions(context.Background(), tt.kv)
				assert.NoError(t, err)

				got, err := s.GetOptions(context.Background())
				assert.NoError(t, err)
				for key, value := range tt.want {
					assert.Equal(t, value, got[key])
				}
			})
		})
	}
}

func TestSetting_GetOptions_Cache(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options"))
		s := New(db)

		// 第一次读取，从数据库加载并缓存
		ret1, err := s.GetOptions(context.Background())
		require.NoError(t, err)
		assert.Equal(t, "無處告別", ret1["site_name"])

		// 第二次读取，应从缓存返回
		ret2, err := s.GetOptions(context.Background())
		require.NoError(t, err)
		assert.Equal(t, "無處告別", ret2["site_name"])

		// 修改返回值不应影响缓存
		ret1["site_name"] = "modified"
		ret3, err := s.GetOptions(context.Background())
		require.NoError(t, err)
		assert.Equal(t, "無處告別", ret3["site_name"])
	})
}

func TestSetting_GetAIConfig(t *testing.T) {
	t.Run("无AI配置返回空值", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options"))
			s := New(db)
			cfg := s.GetAIConfig(context.Background())
			assert.NotNil(t, cfg)
			assert.Empty(t, cfg.Token)
			assert.Empty(t, cfg.Endpoint)
			assert.Empty(t, cfg.Model)
		})
	})

	t.Run("有AI配置返回正确值", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options"))
			s := New(db)

			// 写入 AI 配置
			_, err := s.UpdateOptions(context.Background(), map[string]string{
				"ai_token":    "sk-test-token",
				"ai_endpoint": "https://api.example.com",
				"ai_model":    "gpt-test",
			})
			require.NoError(t, err)

			cfg := s.GetAIConfig(context.Background())
			assert.NotNil(t, cfg)
			assert.Equal(t, "sk-test-token", cfg.Token)
			assert.Equal(t, "https://api.example.com", cfg.Endpoint)
			assert.Equal(t, "gpt-test", cfg.Model)
		})
	})
}
