package mcptool

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"app/config"
	"app/pkg/errors"
	"app/service/remind"
	"app/store"
	"app/store/model"

	"github.com/openai/openai-go/v3"
)

type SmartRemindRule struct {
	Type    int    `json:"type"`
	Month   int    `json:"month"`
	Week    int    `json:"week"`
	Day     int    `json:"day"`
	Hour    int    `json:"hour"`
	Minute  int    `json:"minute"`
	Content string `json:"content"`
}

func SmartCreateRemind(ctx context.Context, conf *config.Config, aiClient openai.Client, s *store.Store, content string) (int64, *errors.Error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return 0, errors.BadRequest("EMPTY_CONTENT", "Content cannot be empty")
	}

	now := time.Now()
	currentDate := now.Format("2006-01-02 15:04")

	prompt := fmt.Sprintf(`你是一个提醒规则解析助手。当前日期时间为：%s。
根据用户输入的提醒内容，解析出以下字段：
1) type：提醒类型，整数：
   - 0 表示固定时间（只提醒一次，使用具体年月日时分）
   - 1 表示每分钟
   - 2 表示每个小时
   - 3 表示每周
   - 4 表示每天
   - 5 表示每月
   - 6 表示每年
2) month：月份（1-12，type 为 5 或 6 或 0 且包含具体月份时使用，否则为 0）
3) week：每周第几天（1=周一 ... 7=周日，当 type=3 时使用，否则为 0）
4) day：每月第几天（1-31，当 type 为 0 或 5 或 6 时根据描述推断，否则为 0）
5) hour：小时（0-23，当 type != 1 时必填；type=1 时可忽略）
6) minute：分钟（0-59，当 type 不为 1 时用于精确到分钟）
7) content：提醒的简短内容，去掉时间描述，只保留要做的事情，例如：
   - 输入：“明天早上9点提醒我购买火车票”，content 应为 “购买火车票”
   - 输入：“每周一早上8点提醒喝水”，content 应为 “喝水”

请只输出一个 JSON 对象，不要输出其他任何文字。例如：
{"type":4,"month":0,"week":0,"day":0,"hour":9,"minute":0,"content":"喝水"}`, currentDate)

	userInput := fmt.Sprintf("提醒内容：%s", content)

	aiReq := openai.ChatCompletionNewParams{
		Model: conf.Common.AIModel,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(prompt),
			openai.UserMessage(userInput),
		},
	}
	if strings.HasPrefix(conf.Common.AIModel, "doubao") {
		aiReq.SetExtraFields(map[string]any{
			"thinking": map[string]any{
				"type": "disabled",
			},
		})
	}

	completion, err := aiClient.Chat.Completions.New(ctx, aiReq)
	if err != nil {
		return 0, errors.InternalServer("AI_REMIND_CREATE_ERROR", err.Error())
	}
	if len(completion.Choices) == 0 {
		return 0, errors.InternalServer("AI_REMIND_CREATE_EMPTY", "no choices returned")
	}

	raw := strings.TrimSpace(completion.Choices[0].Message.Content)
	ruleJSON, ok := extractFirstJSONObject(raw)
	if !ok {
		return 0, errors.InternalServer("AI_REMIND_PARSE_ERROR", "failed to extract json object from ai response")
	}

	var rule SmartRemindRule
	if err := json.Unmarshal([]byte(ruleJSON), &rule); err != nil {
		return 0, errors.InternalServer("AI_REMIND_PARSE_ERROR", "failed to parse ai response")
	}

	if rule.Hour < 0 || rule.Hour > 23 || rule.Minute < 0 || rule.Minute > 59 {
		return 0, errors.BadRequest("INVALID_RULE", "invalid hour or minute")
	}

	finalContent := strings.TrimSpace(rule.Content)
	if finalContent == "" {
		finalContent = content
	}

	insert := &model.Remind{
		Type:      rule.Type,
		Content:   finalContent,
		Month:     rule.Month,
		Week:      rule.Week,
		Day:       rule.Day,
		Hour:      rule.Hour,
		Minute:    rule.Minute,
		Status:    1,
		CreatedAt: time.Now(),
	}
	insert.NextTime = remind.NextTimeFromRule(insert.CreatedAt, insert)

	lastID, err := s.CreateRemind(ctx, insert)
	if err != nil {
		return 0, errors.InternalServer("REMIND_CREATE_ERROR", err.Error())
	}

	return lastID, nil
}

func extractFirstJSONObject(text string) (string, bool) {
	text = strings.TrimSpace(text)
	if text == "" {
		return "", false
	}
	if strings.HasPrefix(text, "{") && strings.HasSuffix(text, "}") {
		return text, true
	}

	start := strings.IndexByte(text, '{')
	end := strings.LastIndexByte(text, '}')
	if start < 0 || end < 0 || end <= start {
		return "", false
	}
	return text[start : end+1], true
}
