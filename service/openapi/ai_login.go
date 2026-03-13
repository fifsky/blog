package openapi

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"app/config"
	"app/pkg/aiutil"
	"app/pkg/errors"
	apiv1 "app/proto/gen/api/v1"
	"app/store"
	"app/store/model"

	"github.com/google/uuid"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

var _ apiv1.AILoginServiceServer = (*AILogin)(nil)

type chatMessage struct {
	Role    string
	Content string
}

type AILogin struct {
	apiv1.UnimplementedAILoginServiceServer
	store     *store.Store
	conf      *config.Config
	client    openai.Client
	cache     *lru.Cache[string, []chatMessage]
	cacheLock sync.Mutex
}

func NewAILogin(s *store.Store, conf *config.Config) *AILogin {
	client := openai.NewClient(
		option.WithAPIKey(conf.Common.AIToken),
		option.WithBaseURL(conf.Common.AIEndpoint),
	)
	cache, _ := lru.New[string, []chatMessage](100)
	return &AILogin{
		store:  s,
		conf:   conf,
		client: client,
		cache:  cache,
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
	prompt := fmt.Sprintf(`你是一个严格的身份验证助手。用户需要通过对话证明自己的身份。

用户身份特征（保密，不可直接透露）:
%s

验证规则：
1. 从身份特征中识别独立的、具体的信息点（如：出生地具体到镇/村、职业具体技术栈、家庭成员的具体情况、具体日期、具体品牌等）
2. 用户必须陈述具体的细节才能算答对，泛泛的陈述不计入
3. 判断标准示例（以下回答只是举例，真是回答需要根据用户身份特征来判断）：
   - 错误："有老婆"、"结婚了" → 这只是泛泛陈述
   - 正确："老婆是湖北人"、"和大学同学结婚" → 陈述了具体细节
   - 错误："有女儿"、"有孩子"
   - 正确："女儿生日是0818"、"女儿08月18日出生" → 陈述了具体细节
   - 错误："有车"、"开车"
   - 正确："开奥迪"、"车是奥迪品牌" → 陈述了具体品牌
   - 错误："去过旅游"
   - 正确："去过新疆"、"最远去了新疆" → 陈述了具体地点
4. 累计答对3条具体信息即验证通过，输出 [VERIFIED]

重要约束：
- 绝对不能直接问用户问题
- 绝对不能透露任何身份特征的具体内容
- 绝对不能复述用户的身份信息
- 只能通过模糊引导让用户多说，如"请继续"、"还有吗"等
- 用户答对时只说"正确"，不能重复他说了什么
- 对泛泛陈述要引导用户说更具体的内容，如"能说得更具体一些吗"

回复格式要求：
- 第一行：已答对的具体信息数量（纯数字，如：0、1、2、3）
- 第二行开始：你的回复内容
- 答对3条时，在最后加上 [VERIFIED]

当前已答对: %d 条
已尝试次数: %d/%d`,
		profile.IdentityDescription,
		int(session.VerifiedScore*10),
		session.AttemptCount,
		profile.MaxAttempts,
	)

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(prompt),
	}

	history := a.getChatHistory(session.SessionId)
	for _, msg := range history {
		if msg.Role == "user" {
			messages = append(messages, openai.UserMessage(msg.Content))
		} else {
			messages = append(messages, openai.AssistantMessage(msg.Content))
		}
	}

	messages = append(messages, openai.UserMessage(userMessage))

	aiReq := openai.ChatCompletionNewParams{
		Model:    a.conf.Common.AIModel,
		Messages: messages,
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
	var correctCount int
	var content string

	if len(lines) >= 2 {
		fmt.Sscanf(lines[0], "%d", &correctCount)
		content = strings.TrimSpace(lines[1])
	} else {
		content = response
	}

	verified := strings.Contains(response, "[VERIFIED]")

	if verified {
		correctCount = 3
	}

	score := float64(correctCount) / 3.0
	if score > 1.0 {
		score = 1.0
	}

	a.addChatMessage(session.SessionId, "user", userMessage)
	a.addChatMessage(session.SessionId, "assistant", content)

	return &EvaluationResult{
		Content:  strings.ReplaceAll(content, "[VERIFIED]", ""),
		Score:    score,
		Verified: verified,
	}, nil
}

func (a *AILogin) generateWelcomeMessage(profile *model.AuthProfile) string {
	return "你好！为了验证你的身份，请告诉我一些关于你自己的事情。"
}

func (a *AILogin) getChatHistory(sessionID string) []chatMessage {
	a.cacheLock.Lock()
	defer a.cacheLock.Unlock()
	if history, ok := a.cache.Get(sessionID); ok {
		return history
	}
	return []chatMessage{}
}

func (a *AILogin) addChatMessage(sessionID string, role, content string) {
	a.cacheLock.Lock()
	defer a.cacheLock.Unlock()
	history, _ := a.cache.Get(sessionID)
	history = append(history, chatMessage{Role: role, Content: content})
	a.cache.Add(sessionID, history)
}

func (a *AILogin) clearChatHistory(sessionID string) {
	a.cacheLock.Lock()
	defer a.cacheLock.Unlock()
	a.cache.Remove(sessionID)
}
