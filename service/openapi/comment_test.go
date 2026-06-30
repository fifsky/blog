package openapi

import (
	"context"
	"testing"

	"app/pkg/dbunit"
	"app/pkg/errors"
	apiv1 "app/proto/gen/api/v1"
	"app/store"
	"app/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComment_List(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("comments", "posts")...)
		svc := NewComment(store.New(db), WithCommentModerator(&MockModerator{ShouldPass: true}))

		resp, err := svc.List(context.Background(), apiv1.CommentListRequest_builder{PostId: 7}.Build())
		require.NoError(t, err)
		assert.Len(t, resp.GetList(), 5)

		// 验证字段：avatar 已生成、email 不出现在响应中
		item := resp.GetList()[0]
		assert.NotEmpty(t, item.GetAvatar())
		assert.Contains(t, item.GetAvatar(), "gravatarproxy")
		// CommentItem 不含 email 字段，确保隐私不泄露

		// 验证回复的回复携带 reply_name
		var found bool
		for _, it := range resp.GetList() {
			if it.GetId() == 12 {
				assert.Equal(t, "站长", it.GetReplyName())
				assert.Equal(t, 4, int(it.GetPid()))
				found = true
			}
		}
		assert.True(t, found, "应能找到回复的回复评论")
	})
}

func TestComment_Create(t *testing.T) {
	tests := []struct {
		name      string
		req       *apiv1.CommentCreateRequest
		moderator ContentModerator
		wantErr   bool
		errCode   string
	}{
		{
			name: "创建主评论成功",
			req: apiv1.CommentCreateRequest_builder{
				PostId:  7,
				Name:    "张三",
				Email:   "zhangsan@example.com",
				Content: "文章写得不错",
			}.Build(),
			moderator: &MockModerator{ShouldPass: true},
			wantErr:   false,
		},
		{
			name: "创建回复成功",
			req: apiv1.CommentCreateRequest_builder{
				PostId:    7,
				Name:      "李四",
				Pid:       4,
				ReplyName: "匿名",
				Content:   "回复匿名",
			}.Build(),
			moderator: &MockModerator{ShouldPass: true},
			wantErr:   false,
		},
		{
			name: "XSS防护",
			req: apiv1.CommentCreateRequest_builder{
				PostId:  7,
				Name:    "<script>alert(1)</script>",
				Content: "<img src=x onerror=alert(1)>",
			}.Build(),
			moderator: &MockModerator{ShouldPass: true},
			wantErr:   false,
		},
		{
			name: "内容审核失败",
			req: apiv1.CommentCreateRequest_builder{
				PostId:  7,
				Name:    "测试",
				Content: "违规内容",
			}.Build(),
			moderator: &MockModerator{ShouldPass: false, Reason: "检测到违规内容"},
			wantErr:   true,
			errCode:   "CONTENT_MODERATION_FAILED",
		},
		{
			name: "内容审核服务异常",
			req: apiv1.CommentCreateRequest_builder{
				PostId:  7,
				Name:    "测试",
				Content: "正常内容",
			}.Build(),
			moderator: &MockModerator{Err: errors.InternalServer("CONTENT_MODERATION_ERROR", "服务异常")},
			wantErr:   true,
			errCode:   "CONTENT_MODERATION_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("comments", "posts")...)
				ctx := context.Background()

				svc := NewComment(store.New(db), WithCommentModerator(tt.moderator))

				beforeResp, err := svc.List(ctx, apiv1.CommentListRequest_builder{PostId: 7}.Build())
				require.NoError(t, err)
				beforeCount := len(beforeResp.GetList())

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
				assert.Greater(t, resp.GetId(), int32(0))

				afterResp, err := svc.List(ctx, apiv1.CommentListRequest_builder{PostId: 7}.Build())
				require.NoError(t, err)
				assert.Equal(t, beforeCount+1, len(afterResp.GetList()))

				latest := afterResp.GetList()[len(afterResp.GetList())-1]
				assert.Equal(t, resp.GetId(), latest.GetId())

				if tt.name == "XSS防护" {
					assert.Equal(t, "&lt;script&gt;alert(1)&lt;/script&gt;", latest.GetName())
					assert.Equal(t, "&lt;img src=x onerror=alert(1)&gt;", latest.GetContent())
				}
			})
		})
	}
}

func TestComment_Create_WebsiteIncludedInModeration(t *testing.T) {
	// 验证网址也被拼入审核内容（通过自定义 mock 捕获传入的 content）
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("comments", "posts")...)
		ctx := context.Background()

		mock := &captureModerator{}
		svc := NewComment(store.New(db), WithCommentModerator(mock))

		_, err := svc.Create(ctx, apiv1.CommentCreateRequest_builder{
			PostId:  7,
			Name:    "张三",
			Website: "https://spam.com",
			Content: "正文",
		}.Build())
		require.NoError(t, err)
		// 昵称、网址、内容都应在审核内容中
		assert.Contains(t, mock.captured, "张三")
		assert.Contains(t, mock.captured, "https://spam.com")
		assert.Contains(t, mock.captured, "正文")
	})
}

func TestComment_New(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("comments", "posts")...)
		svc := NewComment(store.New(db), WithCommentModerator(&MockModerator{ShouldPass: true}))

		resp, err := svc.New(context.Background(), nil)
		require.NoError(t, err)
		assert.NotEmpty(t, resp.GetList())

		// 最新评论按时间倒序，第一条应是 id=10
		first := resp.GetList()[0]
		assert.Equal(t, int32(10), first.GetId())
		assert.NotEmpty(t, first.GetPostTitle())
		assert.NotEmpty(t, first.GetAvatar())
	})
}

// captureModerator 捕获审核内容的 mock，用于验证拼接逻辑
type captureModerator struct {
	captured string
}

func (m *captureModerator) Moderate(_ context.Context, content string) error {
	m.captured = content
	return nil
}
