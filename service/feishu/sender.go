package feishu

import (
	"app/config"
	"context"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// Sender 通过飞书机器人发送卡片消息，只负责发送 JSON 字符串
type Sender struct {
	client *lark.Client
	userID string
}

// NewSender 创建飞书发送器，飞书配置缺失时返回 nil
func NewSender(conf config.FeishuConf) *Sender {
	if conf.Appid == "" || conf.AppSecret == "" || conf.UserID == "" {
		return nil
	}
	return &Sender{
		client: lark.NewClient(conf.Appid, conf.AppSecret),
		userID: conf.UserID,
	}
}

// Send 发送卡片 JSON 到飞书
func (f *Sender) Send(ctx context.Context, cardJSON string) error {
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType("open_id").
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(f.userID).
			MsgType("interactive").
			Content(cardJSON).
			Build()).
		Build()

	resp, err := f.client.Im.Message.Create(ctx, req)
	if err != nil {
		return err
	}
	if !resp.Success() {
		return &Error{Code: resp.Code, Msg: resp.Msg}
	}
	return nil
}

// Error 飞书 API 错误
type Error struct {
	Code int
	Msg  string
}

func (e *Error) Error() string { return e.Msg }
