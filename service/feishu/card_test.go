package feishu

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestRemindCard_BuildCard(t *testing.T) {
	card := NewRemindCard(nil, nil, nil)
	msg := RemindMessage{
		Content: "喝水时间到了",
		Time:    "2026-07-01 12:00",
		Token:   "abc123",
	}

	cardJSON := card.BuildCard(msg)

	var data map[string]any
	if err := json.Unmarshal([]byte(cardJSON), &data); err != nil {
		t.Fatalf("BuildCard 生成的 JSON 无效: %v\n%s", err, cardJSON)
	}
	if data["schema"] != "2.0" {
		t.Errorf("schema 应为 2.0, got %v", data["schema"])
	}
	if !strings.Contains(cardJSON, "喝水时间到了") {
		t.Error("卡片内容应包含提醒内容")
	}
	if !strings.Contains(cardJSON, "remind_completed") {
		t.Error("卡片应包含标记完成按钮")
	}
	if !strings.Contains(cardJSON, "remind_later") {
		t.Error("卡片应包含稍后提醒按钮")
	}
}

func TestRemindCard_BuildResultCard(t *testing.T) {
	card := NewRemindCard(nil, nil, nil)
	msg := RemindMessage{
		Content: "喝水时间到了",
		Time:    "2026-07-01 12:00",
		Result:  "已确认收到提醒",
	}

	cardJSON := card.BuildResultCard(msg)

	var data map[string]any
	if err := json.Unmarshal([]byte(cardJSON), &data); err != nil {
		t.Fatalf("BuildResultCard 生成的 JSON 无效: %v\n%s", err, cardJSON)
	}
	if !strings.Contains(cardJSON, "已确认收到提醒") {
		t.Error("卡片内容应包含结果文本")
	}
	if strings.Contains(cardJSON, "remind_completed") {
		t.Error("结果卡片不应包含按钮")
	}
}

func TestRemindCard_SpecialChars(t *testing.T) {
	card := NewRemindCard(nil, nil, nil)
	msg := RemindMessage{
		Content: `包含"引号"和\反斜杠`,
		Time:    "2026-07-01 12:00",
		Token:   "tok\"en",
	}

	cardJSON := card.BuildCard(msg)

	var data map[string]any
	if err := json.Unmarshal([]byte(cardJSON), &data); err != nil {
		t.Fatalf("特殊字符导致 JSON 无效: %v\n%s", err, cardJSON)
	}
}

func TestLinkCard_BuildCard(t *testing.T) {
	card := NewLinkCard(nil, nil)
	msg := LinkMessage{
		Content: "**站点名称**: 测试站点\n**站点地址**: [example.com](https://example.com)",
		Token:   "link_token_123",
	}

	cardJSON := card.BuildCard(msg)

	var data map[string]any
	if err := json.Unmarshal([]byte(cardJSON), &data); err != nil {
		t.Fatalf("BuildCard(link) 生成的 JSON 无效: %v\n%s", err, cardJSON)
	}
	if !strings.Contains(cardJSON, "测试站点") {
		t.Error("卡片内容应包含站点名称")
	}
	if !strings.Contains(cardJSON, "link_approve") {
		t.Error("卡片应包含通过按钮")
	}
	if !strings.Contains(cardJSON, "link_reject") {
		t.Error("卡片应包含拒绝按钮")
	}
}

func TestLinkCard_BuildResultCard(t *testing.T) {
	card := NewLinkCard(nil, nil)
	msg := LinkMessage{
		Content: "**站点名称**: 测试站点",
		Result:  "已通过审核",
	}

	cardJSON := card.BuildResultCard(msg)

	var data map[string]any
	if err := json.Unmarshal([]byte(cardJSON), &data); err != nil {
		t.Fatalf("BuildResultCard(link) 生成的 JSON 无效: %v\n%s", err, cardJSON)
	}
	if !strings.Contains(cardJSON, "已通过审核") {
		t.Error("卡片内容应包含结果文本")
	}
	if strings.Contains(cardJSON, "link_approve") {
		t.Error("结果卡片不应包含按钮")
	}
}

func TestCardRegistry_Handle(t *testing.T) {
	registry := NewCardRegistry()
	registry.Register(NewRemindCard(nil, nil, nil))
	registry.Register(NewLinkCard(nil, nil))

	// 未注册的 action 返回空
	cardJSON, _, _ := registry.Handle(context.TODO(), "unknown_action", "token")
	if cardJSON != "" {
		t.Error("未注册的 action 应返回空字符串")
	}
}
