package openapi

import (
	"context"
	"testing"

	apiv1 "app/proto/gen/api/v1"
	"app/store"
	"app/testutil"

	"github.com/goapt/dbunit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGuestbook_List(t *testing.T) {
	tests := []struct {
		name      string
		page      int32
		wantCount int
		wantTotal int32
		wantErr   bool
	}{
		{
			name:      "第一页",
			page:      1,
			wantCount: 3,
			wantTotal: 3,
			wantErr:   false,
		},
		{
			name:      "第二页-无数据",
			page:      2,
			wantCount: 0,
			wantTotal: 3,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("guestbook")...)
				svc := NewGuestbook(store.New(db))

				resp, err := svc.List(context.Background(), &apiv1.GuestbookListRequest{Page: tt.page})
				if tt.wantErr {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				assert.Len(t, resp.List, tt.wantCount)
				assert.Equal(t, tt.wantTotal, resp.Total)

				// 验证数据字段
				if len(resp.List) > 0 {
					item := resp.List[0]
					assert.NotEmpty(t, item.Name)
					assert.NotEmpty(t, item.Content)
					assert.NotEmpty(t, item.CreatedAt)
				}
			})
		})
	}
}

func TestGuestbook_Create(t *testing.T) {
	tests := []struct {
		name    string
		req     *apiv1.GuestbookCreateRequest
		wantErr bool
	}{
		{
			name: "创建成功",
			req: &apiv1.GuestbookCreateRequest{
				Name:    "张三",
				Content: "这是测试留言",
			},
			wantErr: false,
		},
		{
			name: "XSS防护",
			req: &apiv1.GuestbookCreateRequest{
				Name:    "<script>alert('xss')</script>",
				Content: "<img src=x onerror=alert('xss')>",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				// 加载schema，并加载guestbook fixtures
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("guestbook")...)
				// 模拟带有 Context Header 的 Context
				ctx := context.Background() // 这里可以根据需要构造 Context，目前 IP 逻辑在中间件中，这里可能取不到 IP 或取到空，暂不影响 XSS 测试

				svc := NewGuestbook(store.New(db))

				// 获取创建前的总数
				beforeResp, err := svc.List(ctx, &apiv1.GuestbookListRequest{Page: 1})
				require.NoError(t, err)
				beforeCount := beforeResp.Total

				resp, err := svc.Create(ctx, tt.req)
				if tt.wantErr {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				assert.Greater(t, resp.Id, int32(0))

				// 验证总数增加
				afterResp, err := svc.List(ctx, &apiv1.GuestbookListRequest{Page: 1})
				require.NoError(t, err)
				assert.Equal(t, beforeCount+1, afterResp.Total)
				
				// 验证刚插入的数据是否被转义
				// 由于 List 返回是按 ID 倒序，所以第一条应该是刚插入的
				latest := afterResp.List[0]
				assert.Equal(t, int32(resp.Id), latest.Id)
				
				if tt.name == "XSS防护" {
					assert.Equal(t, "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;", latest.Name)
					assert.Equal(t, "&lt;img src=x onerror=alert(&#39;xss&#39;)&gt;", latest.Content)
				} else {
					assert.Equal(t, tt.req.Name, latest.Name)
					assert.Equal(t, tt.req.Content, latest.Content)
				}
			})
		})
	}
}

