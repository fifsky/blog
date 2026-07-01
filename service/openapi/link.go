package openapi

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"app/config"
	"app/pkg/aesutil"
	apiv1 "app/proto/gen/api/v1"
	"app/service/feishu"
	"app/store"
	"app/store/model"

	"github.com/goapt/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ apiv1.LinkServiceHTTPServer = (*Link)(nil)

type Link struct {
	store    *store.Store
	conf     *config.Config
	linkCard *feishu.LinkCard
	sender   *feishu.FeishuSender
}

func NewLink(s *store.Store, conf *config.Config, sender *feishu.FeishuSender) *Link {
	return &Link{
		store:    s,
		conf:     conf,
		linkCard: feishu.NewLinkCard(s, conf),
		sender:   sender,
	}
}

// All 获取所有已审核通过的链接
func (l *Link) All(ctx context.Context, _ *emptypb.Empty) (*apiv1.LinkMenuResponse, error) {
	links, err := l.store.GetApprovedLinks(ctx)
	if err != nil {
		return nil, err
	}
	resp := apiv1.LinkMenuResponse_builder{}.Build()
	for _, v := range links {
		resp.SetList(append(resp.GetList(), apiv1.LinkMenuItem_builder{Url: v.Url,
			Content: v.Name}.Build(),
		))
	}
	return resp, nil
}

// Submit 提交友情链接申请（状态为审核中），并发送飞书审核通知
func (l *Link) Submit(ctx context.Context, req *apiv1.LinkSubmitRequest) (*emptypb.Empty, error) {
	m := &model.Link{
		Name:      req.GetName(),
		Url:       req.GetUrl(),
		Desc:      req.GetDesc(),
		Status:    model.LinkStatusPending,
		CreatedAt: time.Now(),
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
	if l.sender == nil {
		return
	}

	token, err := aesutil.AesEncode(l.conf.Common.TokenSecret, strconv.FormatInt(id, 10))
	if err != nil {
		logger.Error("link notify aes encode error", slog.String("err", err.Error()))
		return
	}

	msg := feishu.LinkMessage{
		Name:  name,
		URL:   url,
		Desc:  desc,
		Token: token,
	}
	cardJSON := l.linkCard.BuildCard(msg)

	if err := l.sender.Send(context.Background(), cardJSON); err != nil {
		logger.Error("link notify send error", slog.String("err", err.Error()))
	}
}
