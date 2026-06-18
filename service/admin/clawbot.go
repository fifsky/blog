package admin

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"app/config"
	"app/pkg/aiagent"
	"app/pkg/clawbot"
	apperrors "app/pkg/errors"
	adminv1 "app/proto/gen/admin/v1"
	"app/store"

	"github.com/openai/openai-go/v3"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	clawBotOptionAccountID = "clawbot_account_id"
	clawBotOptionBotToken  = "clawbot_bot_token"
	clawBotTokenChunkKey   = "clawbot_bot_token_"
	clawBotOptionBaseURL   = "clawbot_base_url"
	clawBotOptionUserID    = "clawbot_user_id"
	clawBotOptionSavedAt   = "clawbot_saved_at"

	clawBotContextTTL       = time.Hour
	clawBotMaxContextLength = 20
	clawBotOptionValueLimit = 180
	clawBotTokenChunkCount  = 32
)

var _ adminv1.ClawBotServiceHTTPServer = (*ClawBot)(nil)

// ClawBotOption 定义 ClawBot 管理服务的配置项。
type ClawBotOption func(*ClawBot)

// WithClawBotBaseURL 设置 iLink API 地址。
func WithClawBotBaseURL(baseURL string) ClawBotOption {
	return func(c *ClawBot) {
		c.baseURL = strings.TrimSpace(baseURL)
	}
}

// WithClawBotHTTPClient 设置 HTTP 客户端。
func WithClawBotHTTPClient(client *http.Client) ClawBotOption {
	return func(c *ClawBot) {
		if client != nil {
			c.httpClient = client
		}
	}
}

// WithClawBotMonitor 控制是否启动消息监听。
func WithClawBotMonitor(enabled bool) ClawBotOption {
	return func(c *ClawBot) {
		c.monitorEnabled = enabled
	}
}

// ClawBot 管理微信扫码登录、连接状态和消息监听。
type ClawBot struct {
	store  *store.Store
	conf   *config.Config
	agent  *aiagent.Agent
	memory *aiagent.Memory

	baseURL    string
	botType    string
	routeTag   string
	httpClient *http.Client

	sessionsMu sync.RWMutex
	sessions   map[string]*clawbot.LoginSession

	monitorMu      sync.Mutex
	monitorEnabled bool
	monitorCancel  context.CancelFunc
	monitoring     bool
	monitorStatus  string
	lastEventAt    time.Time
	lastError      string
}

func NewClawBot(s *store.Store, conf *config.Config, agent *aiagent.Agent, opts ...ClawBotOption) *ClawBot {
	if conf == nil {
		conf = &config.Config{}
	}

	c := &ClawBot{
		store:          s,
		conf:           conf,
		agent:          agent,
		memory:         aiagent.NewMemory(clawBotContextTTL, clawBotMaxContextLength),
		botType:        clawbot.DefaultBotType,
		httpClient:     &http.Client{},
		sessions:       make(map[string]*clawbot.LoginSession),
		monitorEnabled: true,
		monitorStatus:  "stopped",
	}
	for _, opt := range opts {
		opt(c)
	}

	if c.monitorEnabled {
		account, err := c.loadAccount(context.Background())
		if err != nil {
			c.setMonitorError(err)
		} else if account != nil {
			c.startMonitor(account)
		}
	}
	return c
}

func (c *ClawBot) StartLogin(ctx context.Context, req *adminv1.ClawBotStartLoginRequest) (*adminv1.ClawBotLoginSession, error) {
	session, err := c.newClient().StartLogin(ctx, req.GetAccountHint())
	if err != nil {
		return nil, err
	}

	c.sessionsMu.Lock()
	c.sessions[session.SessionKey] = session
	c.sessionsMu.Unlock()

	return c.sessionToProto(session, "wait"), nil
}

func (c *ClawBot) CheckLogin(ctx context.Context, req *adminv1.ClawBotCheckLoginRequest) (*adminv1.ClawBotCheckLoginResponse, error) {
	sessionKey := strings.TrimSpace(req.GetSessionKey())
	if sessionKey == "" {
		return nil, apperrors.BadRequest("CLAWBOT_SESSION_REQUIRED", "登录会话不能为空")
	}

	session := c.getSession(sessionKey)
	if session == nil {
		return nil, apperrors.BadRequest("CLAWBOT_SESSION_NOT_FOUND", "登录会话不存在或已过期")
	}

	result, err := c.newClient().CheckLogin(ctx, session)
	if err != nil {
		return nil, err
	}

	builder := adminv1.ClawBotCheckLoginResponse_builder{
		Status:  result.Status,
		Session: c.sessionToProto(session, result.Status),
	}
	if result.Account != nil {
		if err := c.saveAccount(ctx, result.Account); err != nil {
			return nil, err
		}
		c.deleteSession(sessionKey)
		c.startMonitor(result.Account)
		builder.Account = c.accountToProto(result.Account)
		builder.Connected = true
	}
	return builder.Build(), nil
}

func (c *ClawBot) Status(ctx context.Context, _ *emptypb.Empty) (*adminv1.ClawBotStatusResponse, error) {
	account, err := c.loadAccount(ctx)
	if err != nil {
		return nil, err
	}
	return c.statusResponse(account), nil
}

func (c *ClawBot) Disconnect(ctx context.Context, _ *emptypb.Empty) (*adminv1.ClawBotStatusResponse, error) {
	c.stopMonitor()
	values := map[string]string{
		clawBotOptionAccountID: "",
		clawBotOptionBotToken:  "",
		clawBotOptionBaseURL:   "",
		clawBotOptionUserID:    "",
		clawBotOptionSavedAt:   "",
	}
	clearLongOption(values, clawBotTokenChunkKey)
	if _, err := c.store.UpdateOptions(ctx, values); err != nil {
		return nil, err
	}
	return c.statusResponse(nil), nil
}

func (c *ClawBot) newClient() *clawbot.Client {
	return clawbot.NewClient(clawbot.Options{
		BaseURL:    c.baseURL,
		BotType:    c.botType,
		RouteTag:   c.routeTag,
		HTTPClient: c.httpClient,
	})
}

func (c *ClawBot) newAPIClient(account *clawbot.Account) *clawbot.APIClient {
	return clawbot.NewAPIClient(clawbot.APIOptions{
		BaseURL:    account.BaseURL,
		Token:      account.BotToken,
		AccountID:  account.AccountID,
		RouteTag:   c.routeTag,
		HTTPClient: c.httpClient,
	})
}

func (c *ClawBot) saveAccount(ctx context.Context, account *clawbot.Account) error {
	if account == nil {
		return apperrors.BadRequest("CLAWBOT_ACCOUNT_EMPTY", "ClawBot 账号不能为空")
	}
	values := map[string]string{
		clawBotOptionAccountID: account.AccountID,
		clawBotOptionBaseURL:   account.BaseURL,
		clawBotOptionUserID:    account.UserID,
		clawBotOptionSavedAt:   account.SavedAt,
	}
	tokenValues, err := splitLongOption(clawBotOptionBotToken, clawBotTokenChunkKey, account.BotToken)
	if err != nil {
		return err
	}
	for key, value := range tokenValues {
		values[key] = value
	}

	_, err = c.store.UpdateOptions(ctx, values)
	return err
}

func (c *ClawBot) loadAccount(ctx context.Context) (*clawbot.Account, error) {
	options, err := c.store.GetOptions(ctx)
	if err != nil {
		return nil, err
	}
	accountID := strings.TrimSpace(options[clawBotOptionAccountID])
	botToken := strings.TrimSpace(options[clawBotOptionBotToken])
	if botToken == "" {
		botToken = joinLongOption(options, clawBotTokenChunkKey)
	}
	if accountID == "" || botToken == "" {
		return nil, nil
	}
	return &clawbot.Account{
		AccountID: accountID,
		BotToken:  botToken,
		BaseURL:   firstOption(options[clawBotOptionBaseURL], clawbot.DefaultBaseURL),
		UserID:    strings.TrimSpace(options[clawBotOptionUserID]),
		SavedAt:   strings.TrimSpace(options[clawBotOptionSavedAt]),
	}, nil
}

func (c *ClawBot) getSession(sessionKey string) *clawbot.LoginSession {
	c.sessionsMu.RLock()
	session := c.sessions[sessionKey]
	c.sessionsMu.RUnlock()
	if session == nil {
		return nil
	}
	if time.Since(session.StartedAt) > clawbot.DefaultQRSessionTTL {
		c.deleteSession(sessionKey)
		return nil
	}
	return session
}

func (c *ClawBot) deleteSession(sessionKey string) {
	c.sessionsMu.Lock()
	delete(c.sessions, sessionKey)
	c.sessionsMu.Unlock()
}

func (c *ClawBot) sessionToProto(session *clawbot.LoginSession, status string) *adminv1.ClawBotLoginSession {
	if session == nil {
		return nil
	}
	return adminv1.ClawBotLoginSession_builder{
		SessionKey: session.SessionKey,
		QrCode:     session.QRCode,
		QrContent:  session.QRContent,
		ExpiresAt:  session.StartedAt.Add(clawbot.DefaultQRSessionTTL).UTC().Format(time.RFC3339),
		Status:     status,
	}.Build()
}

func (c *ClawBot) accountToProto(account *clawbot.Account) *adminv1.ClawBotAccount {
	if account == nil {
		return nil
	}
	return adminv1.ClawBotAccount_builder{
		AccountId: account.AccountID,
		UserId:    account.UserID,
		BaseUrl:   account.BaseURL,
		SavedAt:   account.SavedAt,
	}.Build()
}

func (c *ClawBot) statusResponse(account *clawbot.Account) *adminv1.ClawBotStatusResponse {
	c.monitorMu.Lock()
	monitoring := c.monitoring
	monitorStatus := c.monitorStatus
	lastEventAt := ""
	if !c.lastEventAt.IsZero() {
		lastEventAt = c.lastEventAt.UTC().Format(time.RFC3339)
	}
	lastError := c.lastError
	c.monitorMu.Unlock()

	if account == nil {
		monitoring = false
		if monitorStatus == "" {
			monitorStatus = "stopped"
		}
	}

	return adminv1.ClawBotStatusResponse_builder{
		Connected:     account != nil,
		Account:       c.accountToProto(account),
		Monitoring:    monitoring,
		MonitorStatus: monitorStatus,
		LastEventAt:   lastEventAt,
		LastError:     lastError,
	}.Build()
}

func (c *ClawBot) startMonitor(account *clawbot.Account) {
	if !c.monitorEnabled || account == nil || c.agent == nil {
		return
	}

	api := c.newAPIClient(account)
	sender := clawbot.NewSender(clawbot.SenderOptions{
		API:       api,
		AccountID: account.AccountID,
		BaseURL:   account.BaseURL,
		Token:     account.BotToken,
	})
	configManager := clawbot.NewConfigManager(api)

	ctx, cancel := context.WithCancel(context.Background())

	c.monitorMu.Lock()
	if c.monitorCancel != nil {
		c.monitorCancel()
	}
	c.monitorCancel = cancel
	c.monitoring = true
	c.monitorStatus = "running"
	c.lastError = ""
	c.monitorMu.Unlock()

	go func() {
		err := clawbot.Listen(ctx, clawbot.ListenOptions{
			API:         api,
			AccountID:   account.AccountID,
			SyncBufPath: c.syncBufPath(account.AccountID),
			OnMessages: func(ctx context.Context, messages []clawbot.WeixinMessage) error {
				c.handleMessages(ctx, account, api, configManager, sender, messages)
				return nil
			},
			OnError: func(err error) {
				c.setMonitorError(err)
			},
			OnStatus: func(lastEventAt time.Time) {
				c.setMonitorStatus("running", lastEventAt)
			},
		})

		c.monitorMu.Lock()
		c.monitoring = false
		c.monitorStatus = "stopped"
		if err != nil && ctx.Err() == nil {
			c.lastError = err.Error()
		}
		c.monitorMu.Unlock()
	}()
}

func (c *ClawBot) stopMonitor() {
	c.monitorMu.Lock()
	if c.monitorCancel != nil {
		c.monitorCancel()
		c.monitorCancel = nil
	}
	c.monitoring = false
	c.monitorStatus = "stopped"
	c.monitorMu.Unlock()
}

func (c *ClawBot) handleMessages(ctx context.Context, account *clawbot.Account, api *clawbot.APIClient, configManager *clawbot.ConfigManager, sender *clawbot.Sender, messages []clawbot.WeixinMessage) {
	for _, message := range messages {
		if err := c.handleMessage(ctx, account, api, configManager, sender, message); err != nil {
			c.setMonitorError(err)
		}
	}
}

func (c *ClawBot) handleMessage(ctx context.Context, account *clawbot.Account, api *clawbot.APIClient, configManager *clawbot.ConfigManager, sender *clawbot.Sender, message clawbot.WeixinMessage) error {
	body := strings.TrimSpace(clawbot.BodyFromItemList(message.ItemList))
	if body == "" {
		return nil
	}

	contextToken := strings.TrimSpace(message.ContextToken)
	if contextToken != "" {
		clawbot.SetContextToken(account.AccountID, message.FromUserID, contextToken)
	} else {
		contextToken = clawbot.GetContextToken(account.AccountID, message.FromUserID)
	}
	if contextToken == "" {
		return fmt.Errorf("clawbot context token missing for user %s", message.FromUserID)
	}

	cancelTyping := c.sendMessageTyping(ctx, api, configManager, message.FromUserID, contextToken)
	defer cancelTyping()

	reply, err := c.runAI(ctx, account.AccountID+":"+message.FromUserID, body)
	if err != nil {
		return err
	}
	reply = clawbot.MarkdownToPlainText(reply)
	if reply == "" {
		return nil
	}

	cancelTyping()
	_, err = sender.Conversation(clawbot.Target{
		ToUserID:     message.FromUserID,
		ContextToken: contextToken,
	}).SendText(ctx, reply)
	return err
}

func (c *ClawBot) sendMessageTyping(ctx context.Context, api *clawbot.APIClient, configManager *clawbot.ConfigManager, userID, contextToken string) func() {
	if api == nil || configManager == nil {
		return func() {}
	}

	config, err := configManager.GetForUser(ctx, userID, contextToken)
	if err != nil {
		c.setMonitorError(fmt.Errorf("clawbot get config: %w", err))
		return func() {}
	}

	typingTicket := strings.TrimSpace(config.TypingTicket)
	if typingTicket == "" {
		return func() {}
	}

	req := clawbot.SendTypingRequest{
		ILinkUserID:  userID,
		TypingTicket: typingTicket,
		Status:       clawbot.TypingStatusTyping,
	}
	if err := api.SendTyping(ctx, req, 0); err != nil {
		c.setMonitorError(fmt.Errorf("clawbot send typing: %w", err))
		return func() {}
	}

	cancelled := false
	return func() {
		if cancelled {
			return
		}
		cancelled = true
		req.Status = clawbot.TypingStatusCancel
		if err := api.SendTyping(ctx, req, 0); err != nil {
			c.setMonitorError(fmt.Errorf("clawbot cancel typing: %w", err))
		}
	}
}

func (c *ClawBot) runAI(ctx context.Context, memoryKey, userMessage string) (string, error) {
	if c.agent == nil {
		return "", fmt.Errorf("clawbot aiagent is nil")
	}

	c.memory.CleanExpired()
	contextMessages := c.memory.Get(memoryKey)
	messages := make([]openai.ChatCompletionMessageParamUnion, 0, len(contextMessages)+1)
	messages = append(messages, contextMessages...)
	messages = append(messages, openai.UserMessage(userMessage))

	var content strings.Builder
	result, err := c.agent.Run(ctx, aiagent.Request{
		SystemPrompt: c.buildSystemPrompt(),
		Messages:     messages,
		UseTools:     true,
	}, aiagent.EventHandler{
		OnContent: func(_ context.Context, delta string) error {
			content.WriteString(delta)
			return nil
		},
	})
	if err != nil {
		return "", err
	}

	finalContent := content.String()
	if finalContent == "" {
		finalContent = result.Content
	}
	if finalContent == "" {
		finalContent = "抱歉，我暂时无法回答这个问题。"
	}

	memoryMessages := []openai.ChatCompletionMessageParamUnion{openai.UserMessage(userMessage)}
	if len(result.Messages) > 0 {
		memoryMessages = append(memoryMessages, result.Messages...)
	} else {
		memoryMessages = append(memoryMessages, openai.AssistantMessage(finalContent))
	}
	c.memory.Append(memoryKey, memoryMessages...)

	return finalContent, nil
}

func (c *ClawBot) buildSystemPrompt() string {
	return fmt.Sprintf(`You are a helpful assistant inside a WeChat conversation. Respond in the same language as the user's message.
Use available tools when current or external information is needed.
Current Time: %s`, time.Now().Format(time.DateTime))
}

func (c *ClawBot) setMonitorError(err error) {
	if err == nil {
		return
	}
	c.monitorMu.Lock()
	c.lastError = err.Error()
	c.monitorMu.Unlock()
}

func (c *ClawBot) setMonitorStatus(status string, lastEventAt time.Time) {
	c.monitorMu.Lock()
	c.monitorStatus = status
	c.lastEventAt = lastEventAt
	c.monitorMu.Unlock()
}

func (c *ClawBot) syncBufPath(accountID string) string {
	storagePath := strings.TrimSpace(c.conf.Common.StoragePath)
	if storagePath == "" {
		storagePath = "storage"
	}
	filename := base64.RawURLEncoding.EncodeToString([]byte(accountID)) + ".syncbuf"
	return filepath.Join(storagePath, "clawbot", filename)
}

func firstOption(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func splitLongOption(baseKey, chunkKey, value string) (map[string]string, error) {
	values := map[string]string{baseKey: ""}
	clearLongOption(values, chunkKey)
	if len(value) <= clawBotOptionValueLimit {
		values[baseKey] = value
		return values, nil
	}
	if len(value) > clawBotOptionValueLimit*clawBotTokenChunkCount {
		return nil, fmt.Errorf("clawbot option value is too long")
	}

	for index, start := 0, 0; start < len(value); index, start = index+1, start+clawBotOptionValueLimit {
		end := min(start+clawBotOptionValueLimit, len(value))
		values[fmt.Sprintf("%s%d", chunkKey, index)] = value[start:end]
	}
	return values, nil
}

func clearLongOption(values map[string]string, chunkKey string) {
	for index := range clawBotTokenChunkCount {
		values[fmt.Sprintf("%s%d", chunkKey, index)] = ""
	}
}

func joinLongOption(options map[string]string, chunkKey string) string {
	var builder strings.Builder
	for index := range clawBotTokenChunkCount {
		part := options[fmt.Sprintf("%s%d", chunkKey, index)]
		if part == "" {
			break
		}
		builder.WriteString(part)
	}
	return builder.String()
}
