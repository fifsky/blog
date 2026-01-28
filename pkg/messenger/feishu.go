package messenger

import (
	"context"
	"encoding/json"

	"app/config"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// FeishuSender sends messages via Feishu bot card message.
type FeishuSender struct {
	client *lark.Client
	userID string // User open_id to receive the message
}

// NewFeishuSender creates a new FeishuSender.
func NewFeishuSender(conf config.FeishuConf) *FeishuSender {
	if conf.Appid == "" || conf.AppSecret == "" || conf.UserID == "" {
		return nil
	}

	client := lark.NewClient(conf.Appid, conf.AppSecret)
	return &FeishuSender{
		client: client,
		userID: conf.UserID,
	}
}

// Send sends the message via Feishu card.
func (f *FeishuSender) Send(ctx context.Context, msg Message) error {
	cardJSON := f.buildCardJSON(msg)

	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType("open_id").
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(f.userID).
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

// buildCardJSON builds the Feishu card JSON.
func (f *FeishuSender) buildCardJSON(msg Message) string {
	// Build action buttons
	columns := make([]map[string]any, 0, len(msg.Actions))
	for _, action := range msg.Actions {
		column := map[string]any{
			"tag":   "column",
			"width": "auto",
			"elements": []map[string]any{
				{
					"tag":    "button",
					"text":   map[string]any{"tag": "plain_text", "content": action.Title},
					"type":   "primary_filled",
					"width":  "fill",
					"size":   "medium",
					"margin": "4px 0px 4px 0px",
					"behaviors": []map[string]any{
						{
							"type":      "open_url",
							"multi_url": map[string]any{"url": action.URL},
						},
					},
				},
			},
			"vertical_spacing": "8px",
			"horizontal_align": "left",
			"vertical_align":   "top",
		}
		// Second button uses default style
		if len(columns) > 0 {
			column["elements"].([]map[string]any)[0]["type"] = "default"
		}
		columns = append(columns, column)
	}

	// Build elements
	elements := []map[string]any{
		{
			"tag":        "markdown",
			"content":    "**" + msg.Title + "**",
			"text_align": "left",
			"text_size":  "heading",
			"margin":     "0px 0px 0px 0px",
		},
	}

	// Add time if provided
	if msg.Time != "" {
		elements = append(elements, map[string]any{
			"tag":        "markdown",
			"content":    msg.Time,
			"text_align": "left",
			"text_size":  "normal",
			"margin":     "0px 0px 0px 0px",
			"icon": map[string]any{
				"tag":   "standard_icon",
				"token": "calendar_outlined",
				"color": "grey",
			},
		})
	}

	// Add content
	elements = append(elements, map[string]any{
		"tag":        "markdown",
		"content":    msg.Content,
		"text_align": "left",
		"text_size":  "normal",
		"margin":     "0px 0px 0px 0px",
	})

	// Add action buttons if any
	if len(columns) > 0 {
		elements = append(elements, map[string]any{
			"tag":                "column_set",
			"flex_mode":          "stretch",
			"horizontal_spacing": "8px",
			"horizontal_align":   "left",
			"columns":            columns,
			"margin":             "0px 0px 0px 0px",
		})
	}

	card := map[string]any{
		"schema": "2.0",
		"config": map[string]any{
			"update_multi": true,
		},
		"body": map[string]any{
			"direction":        "vertical",
			"vertical_spacing": "12px",
			"padding":          "20px 20px 20px 20px",
			"elements":         elements,
		},
		"header": map[string]any{
			"title":    map[string]any{"tag": "plain_text", "content": "个人事项提醒"},
			"subtitle": map[string]any{"tag": "plain_text", "content": ""},
			"template": "yellow",
			"icon":     map[string]any{"tag": "standard_icon", "token": "reminder_outlined"},
			"padding":  "12px 12px 12px 0px",
		},
	}

	jsonBytes, _ := json.Marshal(card)
	return string(jsonBytes)
}
