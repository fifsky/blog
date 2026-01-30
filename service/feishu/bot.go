package feishu

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"app/config"
	"app/pkg/aesutil"
	apiv1 "app/proto/gen/api/v1"
	"app/service/openapi"
	"app/store"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher/callback"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/larksuite/oapi-sdk-go/v3/ws"
	"github.com/samber/lo"
)

// Bot represents the Feishu bot service that listens for messages via WebSocket
// and responds with AI-generated content using streaming card updates.
type Bot struct {
	conf       *config.Config
	larkClient *lark.Client
	wsClient   *ws.Client
	aiChat     *AIChat
	remind     *openapi.Remind
	store      *store.Store
}

// NewBot creates a new Feishu bot instance.
func NewBot(conf *config.Config, s *store.Store) *Bot {
	// Create Lark client for API calls
	larkClient := lark.NewClient(
		conf.Feishu.Appid,
		conf.Feishu.AppSecret,
		lark.WithLogLevel(larkcore.LogLevelInfo),
	)

	// Create AI chat handler
	aiChat := NewAIChat(conf, larkClient)

	// Create remind service for card callback handling
	remind := openapi.NewRemind(s, conf)

	// Create bot instance first so we can reference it in the handler
	bot := &Bot{
		conf:       conf,
		larkClient: larkClient,
		aiChat:     aiChat,
		remind:     remind,
		store:      s,
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
func (b *Bot) handleCardAction(ctx context.Context, event *callback.CardActionTriggerEvent) (*callback.CardActionTriggerResponse, error) {
	// fmt.Printf("[Feishu Bot] Card action received: %s\n", larkcore.Prettify(event))
	// Parse action value from event
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

	// Execute action based on key
	var result *apiv1.TextResponse
	var err error

	req := &apiv1.RemindActionRequest{Token: token}
	switch actionKey {
	case "remind_completed":
		result, err = b.remind.Change(ctx, req)
	case "remind_later":
		result, err = b.remind.Delay(ctx, req)
	default:
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("操作失败: %w", err)
	}

	id, err := aesutil.AesDecode(b.conf.Common.TokenSecret, req.Token)
	if err != nil {
		return nil, fmt.Errorf("token错误:%w", err)
	}

	remind, err := b.store.GetRemind(ctx, lo.Must(strconv.Atoi(id)))
	if err != nil {
		return nil, fmt.Errorf("记录未找到:%w", err)
	}

	// Get response text
	responseText := "操作完成"
	if result != nil {
		responseText = result.Text
	}

	// Return only toast response
	return &callback.CardActionTriggerResponse{
		Toast: &callback.Toast{
			Type:    "success",
			Content: responseText,
		},
		Card: &callback.Card{
			Type: "template",
			Data: &callback.TemplateCard{
				TemplateID: b.conf.Feishu.RemindResultTemplateID,
				TemplateVariable: map[string]any{
					"remind_content": remind.Content,
					"remind_time":    remind.NextTime.Format("2006-01-02 15:04"),
					"remind_result":  responseText,
				},
			},
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
