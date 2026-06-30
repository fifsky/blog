package openapi

import (
	"context"
	"encoding/json"
	"strings"

	"app/pkg/errors"
	"app/store"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// AIModerator 基于 AI 的内容审核器
type AIModerator struct {
	store *store.Store
}

// NewAIModerator 创建 AI 内容审核器
func NewAIModerator(s *store.Store) *AIModerator {
	return &AIModerator{
		store: s,
	}
}

// Moderate 使用 AI 对内容进行审核
func (m *AIModerator) Moderate(ctx context.Context, content string) error {
	if strings.TrimSpace(content) == "" {
		return nil
	}

	aiCfg := m.store.GetAIConfig(ctx)
	if aiCfg.Token == "" {
		return nil // 没有配置 AI token 时跳过审核
	}
	aiClient := openai.NewClient(
		option.WithAPIKey(aiCfg.Token),
		option.WithBaseURL(aiCfg.Endpoint),
	)

	prompt := `你是一个内容安全审核助手。请审核以下用户提交的留言内容，判断是否包含以下任何一种违规内容：
1. 色情、性感内容
2. 涉政敏感内容
3. 暴力恐怖内容
4. 违禁品相关内容
5. 宗教极端内容
6. 引流广告、垃圾广告
7. 辱骂、歧视、仇恨言论
8. 其他不良内容

请只回复 JSON 格式：
- 如果内容合规，回复：{"pass": true}
- 如果内容违规，回复：{"pass": false, "reason": "违规原因简述"}

不要输出任何其他内容。`

	aiReq := openai.ChatCompletionNewParams{
		Model: aiCfg.Model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(prompt),
			openai.UserMessage(content),
		},
	}
	aiReq.SetExtraFields(map[string]any{
		"thinking": map[string]any{
			"type": "disabled",
		},
	})

	completion, err := aiClient.Chat.Completions.New(ctx, aiReq)
	if err != nil {
		return errors.InternalServer("CONTENT_MODERATION_ERROR", "内容审核服务异常")
	}

	if len(completion.Choices) == 0 {
		return errors.InternalServer("CONTENT_MODERATION_ERROR", "内容审核服务无响应")
	}

	result := strings.TrimSpace(completion.Choices[0].Message.Content)

	// 解析 AI 返回的审核结果
	var moderationResult struct {
		Pass   bool   `json:"pass"`
		Reason string `json:"reason"`
	}

	if err := json.Unmarshal([]byte(result), &moderationResult); err != nil {
		return errors.InternalServer("CONTENT_MODERATION_ERROR", "内容审核结果解析失败")
	}

	if !moderationResult.Pass {
		reason := moderationResult.Reason
		if reason == "" {
			reason = "内容包含违规信息"
		}
		return errors.BadRequest("CONTENT_MODERATION_FAILED", reason)
	}

	return nil
}
