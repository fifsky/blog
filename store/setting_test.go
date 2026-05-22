package store

import (
	"context"
	"testing"

	"app/testutil"

	"github.com/goapt/dbunit"
	"github.com/stretchr/testify/assert"
)

func TestSetting_GetOptions(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options")...)
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
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options")...)
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
