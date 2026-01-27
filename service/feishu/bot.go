package feishu

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"app/config"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
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
}

// NewBot creates a new Feishu bot instance.
func NewBot(conf *config.Config) *Bot {
	// Create Lark client for API calls
	larkClient := lark.NewClient(
		conf.Feishu.Appid,
		conf.Feishu.AppSecret,
		lark.WithLogLevel(larkcore.LogLevelInfo),
	)

	// Create AI chat handler
	aiChat := NewAIChat(conf, larkClient)

	// Create bot instance first so we can reference it in the handler
	bot := &Bot{
		conf:       conf,
		larkClient: larkClient,
		aiChat:     aiChat,
	}

	// Create event dispatcher with fluent handler registration (official pattern)
	eventHandler := dispatcher.NewEventDispatcher("", "").
		OnP2MessageReceiveV1(bot.handleMessage)

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
func (b *Bot) Start(ctx context.Context) error {
	fmt.Println("[Feishu Bot] Starting WebSocket connection...")
	return b.wsClient.Start(ctx)
}

// handleMessage handles incoming P2P and group messages.
func (b *Bot) handleMessage(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	// Skip non-text messages
	msgType := event.Event.Message.MessageType
	if msgType == nil || *msgType != "text" {
		return nil
	}

	// Extract message content
	content := event.Event.Message.Content
	if content == nil {
		return nil
	}

	// Parse text content
	var textContent struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal([]byte(*content), &textContent); err != nil {
		fmt.Printf("[Feishu Bot] Failed to parse message content: %v\n", err)
		return nil
	}

	text := strings.TrimSpace(textContent.Text)
	if text == "" {
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

	fmt.Printf("[Feishu Bot] Received message from %s: %s\n", senderID, text)

	// Handle AI chat in a goroutine to not block the event handler
	go func() {
		if err := b.aiChat.HandleMessage(context.Background(), messageID, text); err != nil {
			fmt.Printf("[Feishu Bot] Failed to handle AI chat: %v\n", err)
		}
	}()

	return nil
}
