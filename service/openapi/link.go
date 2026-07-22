package openapi

import (
	"context"
	"log/slog"
	"time"

	apiv1 "app/proto/gen/api/v1"
	"app/service/feishu"
	"app/store"
	"app/store/model"

	"github.com/goapt/logger"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ apiv1.LinkServiceHTTPServer = (*Link)(nil)

type Link struct {
	store    *store.Store
	linkCard *feishu.LinkCard
}

func NewLink(s *store.Store, feishuConf feishu.Config) *Link {
	return &Link{
		store:    s,
		linkCard: feishu.NewLinkCard(s, feishuConf),
	}
}

// All 获取所有已审核通过的链接
func (l *Link) All(ctx context.Context, _ *emptypb.Empty) (*apiv1.LinkMenuResponse, error) {
	links, err := l.store.GetApprovedLinks(ctx)
	if err != nil {
		return nil, err
	}
	items := lo.Map(links, func(v *model.Link, _ int) *apiv1.LinkMenuItem {
		return apiv1.LinkMenuItem_builder{Url: v.Url, Content: v.Name}.Build()
	})
	return apiv1.LinkMenuResponse_builder{List: items}.Build(), nil
}

// Submit 提交友情链接申请（状态为审核中），并发送飞书审核通知
func (l *Link) Submit(ctx context.Context, req *apiv1.LinkSubmitRequest) (*emptypb.Empty, error) {
	m := &model.Link{
		Name:      req.GetName(),
		Url:       req.GetUrl(),
		Desc:      req.GetDesc(),
		Status:    model.LinkStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	id, err := l.store.CreateLink(ctx, m)
	if err != nil {
		return nil, err
	}

	// 发送飞书审核通知
	l.notifyLinkSubmit(id, req.GetName(), req.GetUrl(), req.GetDesc())

	return &emptypb.Empty{}, nil
}

// notifyLinkSubmit 发送友情链接提交审核通知
func (l *Link) notifyLinkSubmit(id int64, name, url, desc string) {
	msg := feishu.LinkMessage{
		Name: name,
		URL:  url,
		Desc: desc,
		ID:   int(id),
	}
	if err := l.linkCard.Send(context.Background(), msg); err != nil {
		logger.Error("link notify send error", slog.String("err", err.Error()))
	}
}
