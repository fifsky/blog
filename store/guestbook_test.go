package store

import (
	"context"
	"testing"
	"time"

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
		keyword   string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "第一页",
			page:      1,
			keyword:   "",
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:      "第二页-无数据",
			page:      2,
			keyword:   "",
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:      "搜索关键字-匹配Name",
			page:      1,
			keyword:   "张三",
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "搜索关键字-匹配Content",
			page:      1,
			keyword:   "留言测试",
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "搜索关键字-无匹配",
			page:      1,
			keyword:   "不存在的关键字",
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("guestbook")...)
				s := New(db)

				got, err := s.ListGuestbook(context.Background(), tt.keyword, tt.page, 10)
				if tt.wantErr {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				assert.Len(t, got, tt.wantCount)
			})
		})
	}
}

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.DateTime, s)
	return t
}

func TestStore_ListGuestbook_Order(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema())
		s := New(db)

		// 清空可能存在的默认数据
		_, err := db.Exec("TRUNCATE TABLE guestbook")
		require.NoError(t, err)

		// 插入测试数据，包含置顶和非置顶数据
		// 注意：created_at 需要不同以验证排序
		g1 := &model.Guestbook{Name: "User1", Content: "Content1", Top: 0, CreatedAt: parseTime("2023-01-01 10:00:00")}
		g2 := &model.Guestbook{Name: "User2", Content: "Content2", Top: 1, CreatedAt: parseTime("2023-01-02 10:00:00")} // 置顶
		g3 := &model.Guestbook{Name: "User3", Content: "Content3", Top: 0, CreatedAt: parseTime("2023-01-03 10:00:00")}

		// 手动插入以确保顺序可控或直接依赖List的排序
		_, _ = s.CreateGuestbook(context.Background(), g1)
		_, _ = s.CreateGuestbook(context.Background(), g2)
		_, _ = s.CreateGuestbook(context.Background(), g3)

		// 更新g2为置顶
		_, err = db.Exec("UPDATE guestbook SET top = 1 WHERE name = 'User2'")
		require.NoError(t, err)

		list, err := s.ListGuestbook(context.Background(), "", 1, 10)
		require.NoError(t, err)
		assert.Len(t, list, 3)

		// 期望顺序：
		// 1. User2 (Top=1)
		// 2. User3 (Top=0, CreatedAt=2023-01-03)
		// 3. User1 (Top=0, CreatedAt=2023-01-01)

		assert.Equal(t, "User2", list[0].Name)
		assert.Equal(t, "User3", list[1].Name)
		assert.Equal(t, "User1", list[2].Name)
	})
}

func TestStore_CountGuestbookTotal(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("guestbook")...)
		s := New(db)

		total, err := s.CountGuestbookTotal(context.Background(), "")
		require.NoError(t, err)
		assert.Equal(t, 3, total)

		total, err = s.CountGuestbookTotal(context.Background(), "张三")
		require.NoError(t, err)
		assert.Equal(t, 1, total)
	})
}

func TestStore_CreateGuestbook(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema())
		s := New(db)

		// 获取创建前的总数
		beforeTotal, err := s.CountGuestbookTotal(context.Background(), "")
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
		afterTotal, err := s.CountGuestbookTotal(context.Background(), "")
		require.NoError(t, err)
		assert.Equal(t, beforeTotal+1, afterTotal)

		// 验证创建的数据
		list, err := s.ListGuestbook(context.Background(), "", 1, 10)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(list), 1)
		// 新创建的数据应该在第一条（按ID降序）
		assert.Equal(t, "测试用户", list[0].Name)
		assert.Equal(t, "这是测试留言内容", list[0].Content)
		assert.Equal(t, "127.0.0.1", list[0].Ip)
	})
}
