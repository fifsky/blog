package messenger

import (
	"context"
	"fmt"

	"app/pkg/wechat"
)

// WechatSender sends messages via WeChat Work robot.
type WechatSender struct {
	robot *wechat.Robot
}

// NewWechatSender creates a new WechatSender.
func NewWechatSender(token string) *WechatSender {
	if token == "" {
		return nil
	}
	return &WechatSender{robot: wechat.NewRobot(token)}
}

// Send sends the message via WeChat Work robot.
func (w *WechatSender) Send(ctx context.Context, msg Message) error {
	actions := make([]map[string]string, 0, len(msg.Actions))
	for _, action := range msg.Actions {
		actions = append(actions, map[string]string{
			"title":     action.Title,
			"actionURL": action.URL,
		})
	}

	content := msg.Content
	if msg.Time != "" {
		content = fmt.Sprintf("提醒时间: %s\n\n%s", msg.Time, content)
	}

	return w.robot.CardMessage(msg.Title, content, actions)
}
