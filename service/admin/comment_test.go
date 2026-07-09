package admin

import (
	"context"
	"testing"

	"app/pkg/dbunit"
	adminv1 "app/proto/gen/admin/v1"
	"app/store"
	"app/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminComment_List(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("comments", "posts"))
		svc := NewComment(store.New(db))

		resp, err := svc.List(context.Background(), adminv1.CommentListRequest_builder{Page: 1}.Build())
		require.NoError(t, err)
		assert.Equal(t, int32(5), resp.GetTotal())
		assert.Len(t, resp.GetList(), 5)

		// 后台响应应包含 email 和 ip（完整信息）
		first := resp.GetList()[0]
		_ = first.GetEmail()
		_ = first.GetIp()
		assert.NotEmpty(t, first.GetPostTitle())

		// 验证搜索（"时光飞逝" 仅在 id=11 的内容中出现）
		resp2, err := svc.List(context.Background(), adminv1.CommentListRequest_builder{Page: 1, Keyword: "时光飞逝"}.Build())
		require.NoError(t, err)
		assert.Equal(t, int32(1), resp2.GetTotal())
	})
}

func TestAdminComment_Delete(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("comments", "posts"))
		svc := NewComment(store.New(db))

		before, err := svc.List(context.Background(), adminv1.CommentListRequest_builder{Page: 1}.Build())
		require.NoError(t, err)
		beforeTotal := before.GetTotal()

		_, err = svc.Delete(context.Background(), adminv1.CommentDeleteRequest_builder{Ids: []int32{4, 9}}.Build())
		require.NoError(t, err)

		after, err := svc.List(context.Background(), adminv1.CommentListRequest_builder{Page: 1}.Build())
		require.NoError(t, err)
		assert.Equal(t, beforeTotal-2, after.GetTotal())
	})
}

func TestAdminComment_DeleteEmptyIds(t *testing.T) {
	// 空 ID 数组调用不报错（store.DeleteComment 对空数组直接返回 nil）
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("comments"))
		svc := NewComment(store.New(db))

		_, err := svc.Delete(context.Background(), adminv1.CommentDeleteRequest_builder{Ids: []int32{}}.Build())
		require.NoError(t, err)
	})
}
