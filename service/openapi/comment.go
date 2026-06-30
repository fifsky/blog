package openapi

import (
	"context"
	"html"
	"strings"
	"time"

	"app/pkg/gravatar"
	apiv1 "app/proto/gen/api/v1"
	"app/server/middleware"
	"app/store"
	"app/store/model"

	"google.golang.org/protobuf/types/known/emptypb"
)

var _ apiv1.CommentServiceHTTPServer = (*Comment)(nil)

// Comment 文章评论服务
type Comment struct {
	store     *store.Store
	moderator ContentModerator
}

// CommentOption 用于配置 Comment 的选项函数
type CommentOption func(*Comment)

// WithCommentModerator 设置自定义的内容审核器（测试用）
func WithCommentModerator(m ContentModerator) CommentOption {
	return func(c *Comment) {
		c.moderator = m
	}
}

// NewComment 创建评论服务，默认使用 AI 审核器
func NewComment(s *store.Store, opts ...CommentOption) *Comment {
	c := &Comment{
		store: s,
	}

	for _, opt := range opts {
		opt(c)
	}

	// 如果没有设置自定义审核器，使用默认的 AI 审核器
	if c.moderator == nil {
		c.moderator = NewAIModerator(s)
	}

	return c
}

// List 获取某篇文章的全部评论
func (c *Comment) List(ctx context.Context, req *apiv1.CommentListRequest) (*apiv1.CommentListResponse, error) {
	comments, err := c.store.ListComments(ctx, int(req.GetPostId()))
	if err != nil {
		return nil, err
	}

	items := make([]*apiv1.CommentItem, 0, len(comments))
	for _, cm := range comments {
		items = append(items, apiv1.CommentItem_builder{
			Id:        int32(cm.Id),
			Pid:       int32(cm.Pid),
			Name:      cm.Name,
			Avatar:    gravatar.AvatarURL(cm.Email, gravatar.DefaultSize),
			Website:   cm.Website,
			Content:   cm.Content,
			ReplyName: cm.ReplyName,
			CreatedAt: cm.CreatedAt.Format(time.DateTime),
		}.Build())
	}

	return apiv1.CommentListResponse_builder{List: items}.Build(), nil
}

// Create 创建评论，提交的昵称/网址/内容拼接后交由 AI 审核
func (c *Comment) Create(ctx context.Context, req *apiv1.CommentCreateRequest) (*apiv1.CommentCreateResponse, error) {
	// 内容审核：昵称、网址、内容拼接后检测
	if c.moderator != nil {
		content := strings.Join([]string{req.GetName(), req.GetContent()}, " ")
		if err := c.moderator.Moderate(ctx, content); err != nil {
			return nil, err
		}
	}

	// 获取客户端 IP
	ip := middleware.ClientIPFromContext(ctx)

	cm := &model.Comment{
		PostId:    int(req.GetPostId()),
		Pid:       int(req.GetPid()),
		Name:      html.EscapeString(req.GetName()),
		Email:     strings.TrimSpace(req.GetEmail()),
		Website:   html.EscapeString(strings.TrimSpace(req.GetWebsite())),
		Content:   html.EscapeString(req.GetContent()),
		ReplyName: html.EscapeString(req.GetReplyName()),
		IP:        ip,
		CreatedAt: time.Now(),
	}

	id, err := c.store.CreateComment(ctx, cm)
	if err != nil {
		return nil, err
	}

	return apiv1.CommentCreateResponse_builder{Id: int32(id)}.Build(), nil
}

// New 获取最新评论（侧边栏用），返回关联文章信息便于跳转
func (c *Comment) New(ctx context.Context, _ *emptypb.Empty) (*apiv1.CommentNewResponse, error) {
	comments, err := c.store.ListNewComments(ctx, 10)
	if err != nil {
		return nil, err
	}

	items := make([]*apiv1.CommentNewItem, 0, len(comments))
	for _, cm := range comments {
		items = append(items, apiv1.CommentNewItem_builder{
			Id:        int32(cm.Id),
			PostId:    int32(cm.PostId),
			Name:      cm.Name,
			Avatar:    gravatar.AvatarURL(cm.Email, gravatar.DefaultSize),
			Content:   cm.Content,
			CreatedAt: cm.CreatedAt.Format(time.DateTime),
			PostTitle: cm.PostTitle,
			PostUrl:   cm.PostUrl,
		}.Build())
	}

	return apiv1.CommentNewResponse_builder{List: items}.Build(), nil
}
