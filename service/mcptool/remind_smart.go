package mcptool

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"app/config"
	"app/pkg/aiutil"
	"app/pkg/errors"
	"app/pkg/promptutil"
	"app/service/remind"
	"app/store"
	"app/store/model"
	_ "embed"

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

//go:embed remind_prompt.md
var remindPrompt string

func SmartCreateRemind(ctx context.Context, conf *config.Config, aiClient openai.Client, s *store.Store, content string) (int64, *errors.Error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return 0, errors.BadRequest("EMPTY_CONTENT", "Content cannot be empty")
	}

	now := time.Now()
	currentDate := now.Format("2006-01-02 15:04")

	prompt := promptutil.ParsePrompt(remindPrompt, map[string]string{
		"current_date": currentDate,
	})

	userInput := fmt.Sprintf("提醒内容：%s", content)

	aiReq := openai.ChatCompletionNewParams{
		Model: conf.Common.AIModel,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(prompt),
			openai.UserMessage(userInput),
		},
	}
	aiutil.ConfigureModelParams(&aiReq, conf.Common.AIModel)

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
