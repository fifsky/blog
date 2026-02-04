package openapi

import (
	"context"
	"testing"

	"app/pkg/errors"
	apiv1 "app/proto/gen/api/v1"
	"app/store"
	"app/testutil"

	"github.com/goapt/dbunit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockModerator 用于测试的模拟审核器
type MockModerator struct {
	ShouldPass bool
	Reason     string
	Err        error
}

func (m *MockModerator) Moderate(ctx context.Context, content string) error {
	if m.Err != nil {
		return m.Err
	}
	if !m.ShouldPass {
		reason := m.Reason
		if reason == "" {
			reason = "内容包含违规信息"
		}
		return errors.BadRequest("CONTENT_MODERATION_FAILED", reason)
	}
	return nil
}

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
				svc := NewGuestbook(store.New(db), nil, WithModerator(&MockModerator{ShouldPass: true}))

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
		name      string
		req       *apiv1.GuestbookCreateRequest
		moderator ContentModerator
		wantErr   bool
		errCode   string
	}{
		{
			name: "创建成功",
			req: &apiv1.GuestbookCreateRequest{
				Name:    "张三",
				Content: "这是测试留言",
			},
			moderator: &MockModerator{ShouldPass: true},
			wantErr:   false,
		},
		{
			name: "XSS防护",
			req: &apiv1.GuestbookCreateRequest{
				Name:    "<script>alert('xss')</script>",
				Content: "<img src=x onerror=alert('xss')>",
			},
			moderator: &MockModerator{ShouldPass: true},
			wantErr:   false,
		},
		{
			name: "内容审核失败",
			req: &apiv1.GuestbookCreateRequest{
				Name:    "测试",
				Content: "违规内容",
			},
			moderator: &MockModerator{ShouldPass: false, Reason: "检测到违规内容"},
			wantErr:   true,
			errCode:   "CONTENT_MODERATION_FAILED",
		},
		{
			name: "内容审核服务异常",
			req: &apiv1.GuestbookCreateRequest{
				Name:    "测试",
				Content: "正常内容",
			},
			moderator: &MockModerator{Err: errors.InternalServer("CONTENT_MODERATION_ERROR", "服务异常")},
			wantErr:   true,
			errCode:   "CONTENT_MODERATION_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				// 加载schema，并加载guestbook fixtures
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("guestbook")...)
				ctx := context.Background()

				svc := NewGuestbook(store.New(db), nil, WithModerator(tt.moderator))

				// 获取创建前的总数
				beforeResp, err := svc.List(ctx, &apiv1.GuestbookListRequest{Page: 1})
				require.NoError(t, err)
				beforeCount := beforeResp.Total

				resp, err := svc.Create(ctx, tt.req)
				if tt.wantErr {
					require.Error(t, err)
					if tt.errCode != "" {
						appErr, ok := err.(*errors.Error)
						require.True(t, ok, "expected *errors.Error type")
						assert.Equal(t, tt.errCode, appErr.Reason)
					}
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

func TestGuestbook_Create_WithoutModerator(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("guestbook")...)
		ctx := context.Background()

		// 不设置审核器，应该也能正常创建
		svc := NewGuestbook(store.New(db), nil)

		resp, err := svc.Create(ctx, &apiv1.GuestbookCreateRequest{
			Name:    "测试",
			Content: "无审核器测试",
		})
		require.NoError(t, err)
		assert.Greater(t, resp.Id, int32(0))
	})
}
