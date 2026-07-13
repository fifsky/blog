package admin

import (
	"context"
	"time"

	adminv1 "app/proto/gen/admin/v1"
	"app/store"
	"app/store/model"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ adminv1.CommentServiceHTTPServer = (*Comment)(nil)

// Comment 评论后台管理服务
type Comment struct {
	store *store.Store
}

// NewComment 创建评论后台服务
func NewComment(s *store.Store) *Comment {
	return &Comment{store: s}
}

// List 后台分页查询评论列表
func (c *Comment) List(ctx context.Context, req *adminv1.CommentListRequest) (*adminv1.CommentListResponse, error) {
	num := 10
	comments, err := c.store.ListAllComments(ctx, req.GetKeyword(), int(req.GetPage()), num)
	if err != nil {
		return nil, err
	}

	items := lo.Map(comments, func(cm model.CommentWithPost, _ int) *adminv1.CommentItem {
		return adminv1.CommentItem_builder{
			Id:        int32(cm.Id),
			PostId:    int32(cm.PostId),
			Pid:       int32(cm.Pid),
			Name:      cm.Name,
			Email:     cm.Email,
			Website:   cm.Website,
			Content:   cm.Content,
			ReplyName: cm.ReplyName,
			Ip:        cm.IP,
			CreatedAt: cm.CreatedAt.Format(time.DateTime),
			PostTitle: cm.PostTitle,
			PostUrl:   cm.PostUrl,
		}.Build()
	})

	total, err := c.store.CountComments(ctx, req.GetKeyword())
	if err != nil {
		return nil, err
	}

	return adminv1.CommentListResponse_builder{List: items, Total: int32(total)}.Build(), nil
}

// Delete 批量删除评论
func (c *Comment) Delete(ctx context.Context, req *adminv1.CommentDeleteRequest) (*emptypb.Empty, error) {
	ids := lo.Map(req.GetIds(), func(id int32, _ int) int { return int(id) })
	if err := c.store.DeleteComment(ctx, ids); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
