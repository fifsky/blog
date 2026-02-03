package store

import (
	"context"
	"testing"

	"app/store/model"
	"app/testutil"

	"github.com/goapt/dbunit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore_ListGuestbook(t *testing.T) {
	tests := []struct {
		name      string
		page      int
		wantCount int
		wantErr   bool
	}{
		{
			name:      "第一页",
			page:      1,
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:      "第二页-无数据",
			page:      2,
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("guestbook")...)
				s := New(db)

				got, err := s.ListGuestbook(context.Background(), tt.page, 10)
				if tt.wantErr {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				assert.Len(t, got, tt.wantCount)

				// 验证数据按ID降序排列
				if len(got) > 1 {
					for i := 0; i < len(got)-1; i++ {
						assert.Greater(t, got[i].Id, got[i+1].Id, "数据应该按ID降序排列")
					}
				}
			})
		})
	}
}

func TestStore_CountGuestbookTotal(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("guestbook")...)
		s := New(db)

		total, err := s.CountGuestbookTotal(context.Background())
		require.NoError(t, err)
		assert.Equal(t, 3, total)
	})
}

func TestStore_CreateGuestbook(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema())
		s := New(db)

		// 获取创建前的总数
		beforeTotal, err := s.CountGuestbookTotal(context.Background())
		require.NoError(t, err)

		gb := &model.Guestbook{
			Name:    "测试用户",
			Content: "这是测试留言内容",
			Ip:      "127.0.0.1",
		}

		id, err := s.CreateGuestbook(context.Background(), gb)
		require.NoError(t, err)
		assert.Greater(t, id, int64(0))

		// 验证总数增加了1
		afterTotal, err := s.CountGuestbookTotal(context.Background())
		require.NoError(t, err)
		assert.Equal(t, beforeTotal+1, afterTotal)

		// 验证创建的数据
		list, err := s.ListGuestbook(context.Background(), 1, 10)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(list), 1)
		// 新创建的数据应该在第一条（按ID降序）
		assert.Equal(t, "测试用户", list[0].Name)
		assert.Equal(t, "这是测试留言内容", list[0].Content)
		assert.Equal(t, "127.0.0.1", list[0].Ip)
	})
}
