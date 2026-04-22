package openapi

import (
	"context"
	"fmt"
	"strings"
	"time"

	"app/config"
	"app/pkg/aiutil"
	"app/pkg/errors"
	apiv1 "app/proto/gen/api/v1"
	"app/store"
	"app/store/model"

	"github.com/google/uuid"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

var _ apiv1.AILoginServiceServer = (*AILogin)(nil)

type AILogin struct {
	apiv1.UnimplementedAILoginServiceServer
	store  *store.Store
	conf   *config.Config
	client openai.Client
}

func NewAILogin(s *store.Store, conf *config.Config) *AILogin {
	client := openai.NewClient(
		option.WithAPIKey(conf.Common.AIToken),
		option.WithBaseURL(conf.Common.AIEndpoint),
	)
	return &AILogin{
		store:  s,
		conf:   conf,
		client: client,
	}
}

func (a *AILogin) Init(ctx context.Context, req *apiv1.AILoginInitRequest) (*apiv1.AILoginInitResponse, error) {
	user, err := a.store.GetUserByName(ctx, req.UserName)
	if err != nil {
		return nil, errors.BadRequest("USER_NOT_FOUND", "用户不存在")
	}
	if user.Status != 1 {
		return nil, errors.BadRequest("USER_DISABLED", "用户已停用")
	}

	profile, err := a.store.GetAuthProfile(ctx, user.Id)
	if err != nil {
		return nil, errors.BadRequest("PROFILE_NOT_SET", "请先设置身份验证特征")
	}

	sessionID := uuid.New().String()
	session := &model.AuthSession{
		SessionId:     sessionID,
		UserId:        user.Id,
		AttemptCount:  0,
		VerifiedScore: 0,
		Status:        model.AuthSessionActive,
		ExpiresAt:     time.Now().Add(5 * time.Minute),
	}
	if _, err := a.store.CreateAuthSession(ctx, session); err != nil {
		return nil, errors.InternalServer("SESSION_CREATE_ERROR", "创建会话失败")
	}

	welcomeMsg := a.generateWelcomeMessage(profile)

	return &apiv1.AILoginInitResponse{
		SessionId:      sessionID,
		WelcomeMessage: welcomeMsg,
	}, nil
}

func (a *AILogin) Chat(ctx context.Context, req *apiv1.AILoginChatRequest) (*apiv1.AILoginChatResponse, error) {
	session, err := a.store.GetAuthSession(ctx, req.SessionId)
	if err != nil {
		return nil, errors.BadRequest("SESSION_INVALID", "会话无效或已过期")
	}

	profile, err := a.store.GetAuthProfile(ctx, session.UserId)
	if err != nil {
		return nil, errors.InternalServer("PROFILE_ERROR", "获取身份特征失败")
	}

	result, err := a.evaluateWithAI(ctx, profile, session, req.Message)
	if err != nil {
		return nil, errors.InternalServer("AI_ERROR", "AI评估失败")
	}

	session.AttemptCount++
	session.VerifiedScore = result.Score

	if result.Verified {
		session.Status = model.AuthSessionSuccess
		a.store.UpdateAuthSession(ctx, session)

		token, _ := signAccessToken(a.conf.Common.TokenSecret, session.UserId)
		return &apiv1.AILoginChatResponse{
			Content:     result.Content,
			Score:       result.Score,
			Verified:    true,
			AccessToken: token,
		}, nil
	}

	if session.AttemptCount >= profile.MaxAttempts {
		session.Status = model.AuthSessionFailed
		a.store.UpdateAuthSession(ctx, session)
		return &apiv1.AILoginChatResponse{
			Content:      "验证失败次数过多，请稍后再试。",
			Failed:       true,
			ErrorMessage: "EXCEED_MAX_ATTEMPTS",
		}, nil
	}

	a.store.UpdateAuthSession(ctx, session)

	return &apiv1.AILoginChatResponse{
		Content: result.Content,
		Score:   result.Score,
	}, nil
}

type EvaluationResult struct {
	Content  string
	Score    float64
	Verified bool
}

func (a *AILogin) evaluateWithAI(ctx context.Context, profile *model.AuthProfile, session *model.AuthSession, userMessage string) (*EvaluationResult, error) {
	prompt := fmt.Sprintf(`你是一个身份验证助手。用户需要通过对话证明自己的身份。

用户身份特征: %s

你的任务:
1. 根据身份特征提出验证问题或评估用户的回答
2. 评估用户回答的可信度，返回一个0-100的分数
3. 当验证得分达到 %.0f%% 时，在回复最后输出 [VERIFIED]
4. 如果用户明显不是本人，在回复最后输出 [REJECTED]

当前验证进度: %.0f%%
已尝试次数: %d/%d

回复格式要求:
- 第一行必须是分数（纯数字，如：85）
- 第二行开始是你的回复内容
- 验证通过时在最后加上 [VERIFIED]
- 验证失败时在最后加上 [REJECTED]`,
		profile.IdentityDescription,
		profile.VerificationThreshold*100,
		session.VerifiedScore*100,
		session.AttemptCount,
		profile.MaxAttempts,
	)

	aiReq := openai.ChatCompletionNewParams{
		Model: a.conf.Common.AIModel,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(prompt),
			openai.UserMessage(userMessage),
		},
	}
	aiutil.ConfigureModelParams(&aiReq, a.conf.Common.AIModel)

	completion, err := a.client.Chat.Completions.New(ctx, aiReq)
	if err != nil {
		return nil, err
	}

	if len(completion.Choices) == 0 {
		return nil, fmt.Errorf("no response")
	}

	response := completion.Choices[0].Message.Content

	lines := strings.SplitN(response, "\n", 2)
	var score float64
	var content string

	if len(lines) >= 2 {
		fmt.Sscanf(lines[0], "%f", &score)
		content = strings.TrimSpace(lines[1])
	} else {
		content = response
	}

	verified := strings.Contains(response, "[VERIFIED]")
	rejected := strings.Contains(response, "[REJECTED]")

	if rejected {
		score = 0
	}

	if verified {
		score = 1.0
	} else {
		score = score / 100.0
	}

	return &EvaluationResult{
		Content:  strings.ReplaceAll(content, "[VERIFIED]", ""),
		Score:    score,
		Verified: verified,
	}, nil
}

func (a *AILogin) generateWelcomeMessage(profile *model.AuthProfile) string {
	return fmt.Sprintf("你好！为了验证你的身份，请告诉我一些关于你自己的事情。\n\n提示：系统已记录了你的身份特征，请通过对话让我相信你就是本人。你有 %d 次尝试机会。", profile.MaxAttempts)
}
