package messenger

import (
	"context"
	"encoding/json"

	"app/config"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher/callback"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// FeishuSender sends messages via Feishu bot card message.
type FeishuSender struct {
	client *lark.Client
	conf   config.FeishuConf
}

// NewFeishuSender creates a new FeishuSender.
func NewFeishuSender(conf config.FeishuConf) *FeishuSender {
	if conf.Appid == "" || conf.AppSecret == "" || conf.UserID == "" {
		return nil
	}

	client := lark.NewClient(conf.Appid, conf.AppSecret)
	return &FeishuSender{
		client: client,
		conf:   conf,
	}
}

// Send sends the message via Feishu card.
func (f *FeishuSender) Send(ctx context.Context, msg Message) error {
	cardJSON := f.buildCardJSON(msg)

	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType("open_id").
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(f.conf.UserID).
			MsgType("interactive").
			Content(cardJSON).
			Build()).
		Build()

	resp, err := f.client.Im.V1.Message.Create(ctx, req)
	if err != nil {
		return err
	}

	if !resp.Success() {
		return &FeishuError{Code: resp.Code, Msg: resp.Msg}
	}

	return nil
}

// FeishuError represents a Feishu API error.
type FeishuError struct {
	Code int
	Msg  string
}

func (e *FeishuError) Error() string {
	return e.Msg
}

// buildCardJSON builds the Feishu card JSON with callback buttons.
func (f *FeishuSender) buildCardJSON(msg Message) string {
	card := &callback.Card{
		Type: "template",
		Data: &callback.TemplateCard{
			TemplateID: f.conf.RemindTemplateID,
			TemplateVariable: map[string]any{
				"token":          msg.Token,
				"remind_content": msg.Content,
				"remind_time":    msg.Time,
			},
		},
	}

	jsonBytes, _ := json.Marshal(card)
	return string(jsonBytes)
}
