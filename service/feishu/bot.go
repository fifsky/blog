package feishu

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"app/config"
	"app/pkg/aiagent"
	"app/store"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher/callback"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/larksuite/oapi-sdk-go/v3/ws"
)

// Bot represents the Feishu bot service that listens for messages via WebSocket
// and responds with AI-generated content using streaming card updates.
type Bot struct {
	conf       *config.Config
	larkClient *lark.Client
	wsClient   *ws.Client
	aiChat     *AIChat
	registry   *CardRegistry
}

// NewBot creates a new Feishu bot instance.
func NewBot(conf *config.Config, s *store.Store, agent *aiagent.Agent, registry *CardRegistry) *Bot {
	// Create Lark client for API calls
	larkClient := lark.NewClient(
		conf.Feishu.Appid,
		conf.Feishu.AppSecret,
		lark.WithLogLevel(larkcore.LogLevelInfo),
	)

	// Create AI chat handler
	aiChat := NewAIChat(agent, larkClient, s)

	// Create bot instance first so we can reference it in the handler
	bot := &Bot{
		conf:       conf,
		larkClient: larkClient,
		aiChat:     aiChat,
		registry:   registry,
	}

	// Create event dispatcher with fluent handler registration
	eventHandler := dispatcher.NewEventDispatcher("", "").
		OnP2MessageReceiveV1(bot.handleMessage).
		OnP2CardActionTrigger(bot.handleCardAction)

	// Create WebSocket client
	bot.wsClient = ws.NewClient(
		conf.Feishu.Appid,
		conf.Feishu.AppSecret,
		ws.WithEventHandler(eventHandler),
		ws.WithLogLevel(larkcore.LogLevelInfo),
	)

	return bot
}

// Start starts the WebSocket connection and begins listening for messages.
// This method blocks until the connection is closed or an error occurs.
func (b *Bot) Start(ctx context.Context) {
	fmt.Println("[Feishu Bot] Starting WebSocket connection...")
	err := b.wsClient.Start(ctx)
	if err != nil {
		fmt.Printf("[Feishu Bot] wsClient.Start failed: %s\n", err.Error())
	}
}

// handleCardAction handles card button callback actions.
// 逻辑固定不变：解析 action/token -> 委托 registry 分发 -> 返回结果卡片。
// 新增卡片类型只需实现 ActionHandler 并在 NewBot 中注册，无需修改此方法。
func (b *Bot) handleCardAction(ctx context.Context, event *callback.CardActionTriggerEvent) (*callback.CardActionTriggerResponse, error) {
	actionValue := event.Event.Action.Value
	if actionValue == nil {
		return nil, nil
	}

	actionKey, ok := actionValue["action"].(string)
	if !ok {
		return nil, nil
	}

	token, ok := actionValue["token"].(string)
	if !ok {
		return nil, nil
	}

	fmt.Printf("[Feishu Bot] Card action: %s, token: %s\n", actionKey, token)

	cardJSON, resultText, err := b.registry.Handle(ctx, actionKey, token)
	if err != nil {
		return nil, err
	}
	if cardJSON == "" {
		return nil, nil
	}

	var cardData map[string]any
	if err := json.Unmarshal([]byte(cardJSON), &cardData); err != nil {
		return nil, fmt.Errorf("解析卡片JSON失败: %w", err)
	}

	return &callback.CardActionTriggerResponse{
		Toast: &callback.Toast{
			Type:    "success",
			Content: resultText,
		},
		Card: &callback.Card{
			Type: "raw",
			Data: cardData,
		},
	}, nil
}

// handleMessage handles incoming P2P and group messages.
func (b *Bot) handleMessage(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	msgType := event.Event.Message.MessageType
	if msgType == nil {
		return nil
	}

	// Get sender info
	senderID := ""
	if event.Event.Sender != nil && event.Event.Sender.SenderId != nil && event.Event.Sender.SenderId.OpenId != nil {
		senderID = *event.Event.Sender.SenderId.OpenId
	}

	messageID := ""
	if event.Event.Message.MessageId != nil {
		messageID = *event.Event.Message.MessageId
	}

	content := event.Event.Message.Content
	if content == nil {
		return nil
	}

	var userMessage string
	var imageBase64 string

	switch *msgType {
	case "text":
		// Parse text content
		var textContent struct {
			Text string `json:"text"`
		}
		if err := json.Unmarshal([]byte(*content), &textContent); err != nil {
			fmt.Printf("[Feishu Bot] Failed to parse text message content: %v\n", err)
			return nil
		}
		userMessage = strings.TrimSpace(textContent.Text)

	case "image":
		// Parse image content to get image_key
		var imageContent struct {
			ImageKey string `json:"image_key"`
		}
		if err := json.Unmarshal([]byte(*content), &imageContent); err != nil {
			fmt.Printf("[Feishu Bot] Failed to parse image message content: %v\n", err)
			return nil
		}

		// Download image and convert to base64
		base64Data, err := b.downloadImageAsBase64(ctx, messageID, imageContent.ImageKey)
		if err != nil {
			fmt.Printf("[Feishu Bot] Failed to download image: %v\n", err)
			return nil
		}
		imageBase64 = base64Data
		userMessage = "[User sent an image, please describe or analyze this image]"

	case "location":
		// Parse location content
		var locationContent struct {
			Name      string `json:"name"`
			Longitude string `json:"longitude"`
			Latitude  string `json:"latitude"`
		}
		if err := json.Unmarshal([]byte(*content), &locationContent); err != nil {
			fmt.Printf("[Feishu Bot] Failed to parse location message content: %v\n", err)
			return nil
		}

		// Build location context message - make it clear this is supplementary info to the conversation
		userMessage = fmt.Sprintf("[User shared a location - this may be context for the previous question or a new request]\nLocation: %s\nLongitude: %s\nLatitude: %s\n\nPlease use this location information in the context of our conversation. If I asked about something location-related before (like weather, restaurants, etc.), please answer based on this location.",
			locationContent.Name,
			locationContent.Longitude,
			locationContent.Latitude,
		)

	default:
		// Unsupported message types
		return nil
	}

	if userMessage == "" && imageBase64 == "" {
		return nil
	}

	fmt.Printf("[Feishu Bot] Received %s message from %s\n", *msgType, senderID)

	// Handle AI chat in a goroutine to not block the event handler
	go func() {
		if err := b.aiChat.HandleMessage(context.Background(), senderID, messageID, userMessage, imageBase64); err != nil {
			fmt.Printf("[Feishu Bot] Failed to handle AI chat: %v\n", err)
		}
	}()

	return nil
}

// downloadImageAsBase64 downloads an image from Feishu and returns it as base64 encoded string
func (b *Bot) downloadImageAsBase64(ctx context.Context, messageID, imageKey string) (string, error) {
	req := larkim.NewGetMessageResourceReqBuilder().
		MessageId(messageID).
		FileKey(imageKey).
		Type("image").
		Build()

	resp, err := b.larkClient.Im.V1.MessageResource.Get(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to get message resource: %w", err)
	}

	if !resp.Success() {
		return "", fmt.Errorf("get message resource failed: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	// Read image data
	imageData, err := io.ReadAll(resp.File)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %w", err)
	}

	// Convert to base64
	base64Data := base64.StdEncoding.EncodeToString(imageData)
	return base64Data, nil
}
