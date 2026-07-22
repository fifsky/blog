package openapi

import (
	"context"
	"fmt"
	"html"
	"log/slog"
	"strings"
	"time"

	"app/config"
	"app/pkg/gravatar"
	apiv1 "app/proto/gen/api/v1"
	"app/server/middleware"
	"app/service/feishu"
	"app/store"
	"app/store/model"

	"github.com/goapt/logger"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ apiv1.CommentServiceHTTPServer = (*Comment)(nil)

// Comment 文章评论服务
type Comment struct {
	store     *store.Store
	conf      *config.Config
	moderator ContentModerator
	card      *feishu.CommentCard
	sender    *feishu.Sender
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
func NewComment(s *store.Store, conf *config.Config, sender *feishu.Sender, opts ...CommentOption) *Comment {
	c := &Comment{
		store:  s,
		conf:   conf,
		card:   feishu.NewCommentCard(),
		sender: sender,
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

	items := lo.Map(comments, func(cm model.Comment, _ int) *apiv1.CommentItem {
		return apiv1.CommentItem_builder{
			Id:        int32(cm.Id),
			Pid:       int32(cm.Pid),
			Name:      cm.Name,
			Avatar:    gravatar.AvatarURL(cm.Email, gravatar.DefaultSize),
			Website:   cm.Website,
			Content:   cm.Content,
			ReplyName: cm.ReplyName,
			CreatedAt: cm.CreatedAt.Format(time.DateTime),
		}.Build()
	})

	return apiv1.CommentListResponse_builder{List: items}.Build(), nil
}

// Create 创建评论，提交的昵称/内容拼接后交由 AI 审核
func (c *Comment) Create(ctx context.Context, req *apiv1.CommentCreateRequest) (*apiv1.CommentCreateResponse, error) {
	// 内容审核：昵称、内容拼接后检测（网址不参与，避免被误判为广告）
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

	// 发送飞书评论通知
	c.notifyComment(ctx, req.GetName(), req.GetContent(), int(req.GetPostId()), cm.CreatedAt)

	return apiv1.CommentCreateResponse_builder{Id: int32(id)}.Build(), nil
}

// notifyComment 发送新评论飞书通知
func (c *Comment) notifyComment(ctx context.Context, name, content string, postID int, createdAt time.Time) {
	if c.sender == nil || c.card == nil {
		return
	}

	// 查询文章信息用于展示
	post, err := c.store.GetPost(ctx, postID, "")
	if err != nil {
		logger.Error("comment notify get article error", slog.String("err", err.Error()))
		return
	}

	msg := feishu.CommentMessage{
		Name:      name,
		Content:   content,
		PostTitle: post.Title,
		PostURL:   fmt.Sprintf("https://fifsky.com/article/%d/%s", post.Id, post.Url),
		Time:      createdAt.Format("2006-01-02 15:04"),
	}
	cardJSON := c.card.BuildCard(msg)

	if err := c.sender.Send(ctx, cardJSON); err != nil {
		logger.Error("comment notify send error", slog.String("err", err.Error()))
	}
}

// New 获取最新评论（侧边栏用），返回关联文章信息便于跳转
func (c *Comment) New(ctx context.Context, _ *emptypb.Empty) (*apiv1.CommentNewResponse, error) {
	comments, err := c.store.ListNewComments(ctx, 10)
	if err != nil {
		return nil, err
	}

	items := lo.Map(comments, func(cm model.CommentWithPost, _ int) *apiv1.CommentNewItem {
		return apiv1.CommentNewItem_builder{
			Id:        int32(cm.Id),
			PostId:    int32(cm.PostId),
			Name:      cm.Name,
			Avatar:    gravatar.AvatarURL(cm.Email, gravatar.DefaultSize),
			Content:   cm.Content,
			CreatedAt: cm.CreatedAt.Format(time.DateTime),
			PostTitle: cm.PostTitle,
			PostUrl:   cm.PostUrl,
		}.Build()
	})

	return apiv1.CommentNewResponse_builder{List: items}.Build(), nil
}
