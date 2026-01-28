package messenger

import (
	"context"
	"fmt"

	"app/pkg/bark"
)

// BarkSender sends messages via Bark push notification.
type BarkSender struct {
	client *bark.Client
}

// NewBarkSender creates a new BarkSender.
func NewBarkSender(client *bark.Client) *BarkSender {
	return &BarkSender{client: client}
}

// Send sends the message via Bark.
func (b *BarkSender) Send(ctx context.Context, msg Message) error {
	// Build markdown content with actions as links
	content := msg.Content
	if msg.Time != "" {
		content = fmt.Sprintf("提醒时间: %s\n\n%s", msg.Time, content)
	}

	if len(msg.Actions) > 0 {
		content += "\n\n"
		for i, action := range msg.Actions {
			if i > 0 {
				content += "  "
			}
			content += fmt.Sprintf("[%s](%s)", action.Title, action.URL)
		}
	}

	barkMsg := bark.Message{
		Title:    msg.Title,
		Badge:    1,
		Markdown: content,
	}

	return b.client.Send(barkMsg)
}
