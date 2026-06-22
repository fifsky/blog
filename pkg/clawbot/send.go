package clawbot

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	reCodeBlock = regexp.MustCompile("(?s)```[^\\n]*\\n?(.*?)```")
	reImageMD   = regexp.MustCompile(`!\[[^\]]*\]\([^)]*\)`)
	reLinkMD    = regexp.MustCompile(`\[([^\]]+)\]\([^)]*\)`)
	reTableSep  = regexp.MustCompile(`(?m)^\|[\s:|-]+\|$`)
	reTableRow  = regexp.MustCompile(`(?m)^\|(.+)\|$`)
	reEmphasis  = regexp.MustCompile(`[*_~` + "`" + `]+`)
)

var contextTokenStore struct {
	sync.RWMutex
	values map[string]string
}

func init() {
	contextTokenStore.values = make(map[string]string)
}

func MarkdownToPlainText(text string) string {
	result := reCodeBlock.ReplaceAllString(text, "$1")
	result = reImageMD.ReplaceAllString(result, "")
	result = reLinkMD.ReplaceAllString(result, "$1")
	result = reTableSep.ReplaceAllString(result, "")
	result = reTableRow.ReplaceAllStringFunc(result, func(row string) string {
		trimmed := strings.TrimPrefix(strings.TrimSuffix(row, "|"), "|")
		cells := strings.Split(trimmed, "|")
		for i := range cells {
			cells[i] = strings.TrimSpace(cells[i])
		}
		return strings.Join(cells, "  ")
	})
	result = reEmphasis.ReplaceAllString(result, "")
	return strings.TrimSpace(result)
}

type messageAPI interface {
	SendMessage(ctx context.Context, req SendMessageRequest, timeout time.Duration) error
	GetUploadURL(ctx context.Context, req GetUploadURLRequest, timeout time.Duration) (*GetUploadURLResponse, error)
}

type senderOptions struct {
	API            messageAPI
	BaseURL        string
	Token          string
	RouteTag       string
	ChannelVersion string
	HTTPClient     *http.Client
	Timeout        time.Duration
	AccountID      string
	CDNBaseURL     string
}

type sender struct {
	api        messageAPI
	httpClient *http.Client
	timeout    time.Duration
	cdnBaseURL string
}

type Target struct {
	ToUserID     string
	ContextToken string
}

type conversation struct {
	sender *sender
	target Target
}

type MessageContext struct {
	Body              string
	From              string
	To                string
	AccountID         string
	OriginatingTo     string
	MessageSID        string
	Timestamp         int64
	Provider          string
	ChatType          string
	SessionKey        string
	ContextToken      string
	MediaURL          string
	MediaPath         string
	MediaType         string
	CommandBody       string
	CommandAuthorized bool
}

type InboundMediaOptions struct {
	DecryptedPicPath   string
	DecryptedVoicePath string
	VoiceMediaType     string
	DecryptedFilePath  string
	FileMediaType      string
	DecryptedVideoPath string
}

func newSender(opts senderOptions) *sender {
	api := opts.API
	if api == nil {
		api = NewClient(Options{
			BaseURL:        opts.BaseURL,
			Token:          opts.Token,
			RouteTag:       opts.RouteTag,
			ChannelVersion: opts.ChannelVersion,
			HTTPClient:     opts.HTTPClient,
			AccountID:      opts.AccountID,
		})
	}

	httpClient := opts.HTTPClient
	if httpClient == nil {
		if client, ok := api.(*Client); ok && client.httpClient != nil {
			httpClient = client.httpClient
		} else {
			httpClient = &http.Client{}
		}
	}

	cdnBaseURL := strings.TrimSpace(opts.CDNBaseURL)
	if cdnBaseURL == "" {
		if client, ok := api.(*Client); ok {
			cdnBaseURL = strings.TrimSpace(client.cdnBaseURL)
		}
	}
	if cdnBaseURL == "" {
		cdnBaseURL = strings.TrimSpace(opts.BaseURL)
	}
	if cdnBaseURL == "" {
		cdnBaseURL = DefaultCDNBaseURL
	}

	return &sender{
		api:        api,
		httpClient: httpClient,
		timeout:    opts.Timeout,
		cdnBaseURL: cdnBaseURL,
	}
}

func (c *Client) SendText(ctx context.Context, target Target, text string) (string, error) {
	return c.newSender().sendText(ctx, target, text)
}

func (c *Client) SendImage(ctx context.Context, target Target, text string, uploaded UploadedFileInfo) (string, error) {
	return c.newSender().sendImage(ctx, target, text, uploaded)
}

func (c *Client) SendVideo(ctx context.Context, target Target, text string, uploaded UploadedFileInfo) (string, error) {
	return c.newSender().sendVideo(ctx, target, text, uploaded)
}

func (c *Client) SendFile(ctx context.Context, target Target, text, fileName string, uploaded UploadedFileInfo) (string, error) {
	return c.newSender().sendFile(ctx, target, text, fileName, uploaded)
}

func (c *Client) SendMediaFile(ctx context.Context, target Target, filePath, text string) (string, error) {
	return c.newSender().sendMediaFile(ctx, target, filePath, text)
}

func (c *Client) StartTyping(ctx context.Context, target Target) (func(), error) {
	if err := target.validate("startTyping"); err != nil {
		return func() {}, err
	}

	config, err := c.getConfigForUser(ctx, target.ToUserID, target.ContextToken)
	if err != nil {
		return func() {}, err
	}
	typingTicket := strings.TrimSpace(config.TypingTicket)
	if typingTicket == "" {
		return func() {}, nil
	}

	req := SendTypingRequest{
		ILinkUserID:  target.ToUserID,
		TypingTicket: typingTicket,
		Status:       TypingStatusTyping,
	}
	if err := c.SendTyping(ctx, req, 0); err != nil {
		return func() {}, err
	}

	cancelled := false
	return func() {
		if cancelled {
			return
		}
		cancelled = true
		req.Status = TypingStatusCancel
		_ = c.SendTyping(ctx, req, 0)
	}, nil
}

func (c *Client) getConfigForUser(ctx context.Context, userID, contextToken string) (cachedConfig, error) {
	if c.configManager == nil {
		c.configManager = newConfigManager(c)
	}
	return c.configManager.GetForUser(ctx, userID, contextToken)
}

func (c *Client) newSender() *sender {
	return newSender(senderOptions{
		API:        c,
		HTTPClient: c.httpClient,
		Timeout:    c.timeout,
		CDNBaseURL: c.cdnBaseURL,
	})
}

func (s *sender) conversation(target Target) *conversation {
	return &conversation{
		sender: s,
		target: target,
	}
}

func (c *conversation) SendText(ctx context.Context, text string) (string, error) {
	return c.sender.sendText(ctx, c.target, text)
}

func (c *conversation) SendImage(ctx context.Context, text string, uploaded UploadedFileInfo) (string, error) {
	return c.sender.sendImage(ctx, c.target, text, uploaded)
}

func (c *conversation) SendVideo(ctx context.Context, text string, uploaded UploadedFileInfo) (string, error) {
	return c.sender.sendVideo(ctx, c.target, text, uploaded)
}

func (c *conversation) SendFile(ctx context.Context, text, fileName string, uploaded UploadedFileInfo) (string, error) {
	return c.sender.sendFile(ctx, c.target, text, fileName, uploaded)
}

func (c *conversation) SendMediaFile(ctx context.Context, filePath, text string) (string, error) {
	return c.sender.sendMediaFile(ctx, c.target, filePath, text)
}

func (s *sender) sendText(ctx context.Context, target Target, text string) (string, error) {
	if err := target.validate("sendText"); err != nil {
		return "", err
	}

	clientID := GenerateID("openclaw-weixin")
	req := buildTextMessageRequest(target, text, clientID)
	if err := s.api.SendMessage(ctx, req, s.timeout); err != nil {
		return "", err
	}
	return clientID, nil
}

func (s *sender) sendImage(ctx context.Context, target Target, text string, uploaded UploadedFileInfo) (string, error) {
	if err := target.validate("sendImage"); err != nil {
		return "", err
	}
	return s.sendMediaItems(ctx, target, text, buildImageMessageItem(uploaded))
}

func (s *sender) sendVideo(ctx context.Context, target Target, text string, uploaded UploadedFileInfo) (string, error) {
	if err := target.validate("sendVideo"); err != nil {
		return "", err
	}
	return s.sendMediaItems(ctx, target, text, buildVideoMessageItem(uploaded))
}

func (s *sender) sendFile(ctx context.Context, target Target, text, fileName string, uploaded UploadedFileInfo) (string, error) {
	if err := target.validate("sendFile"); err != nil {
		return "", err
	}
	return s.sendMediaItems(ctx, target, text, buildFileMessageItem(fileName, uploaded))
}

func (s *sender) sendMediaFile(ctx context.Context, target Target, filePath, text string) (string, error) {
	if err := target.validate("sendMediaFile"); err != nil {
		return "", err
	}

	mime := MIMEFromFilename(filePath)
	switch {
	case strings.HasPrefix(mime, "video/"):
		uploaded, err := uploadMediaToCDNWithAPI(ctx, filePath, target.ToUserID, s.cdnBaseURL, UploadMediaTypeVideo, s.api, s.httpClient, s.timeout)
		if err != nil {
			return "", err
		}
		return s.sendVideo(ctx, target, text, *uploaded)
	case strings.HasPrefix(mime, "image/"):
		uploaded, err := uploadMediaToCDNWithAPI(ctx, filePath, target.ToUserID, s.cdnBaseURL, UploadMediaTypeImage, s.api, s.httpClient, s.timeout)
		if err != nil {
			return "", err
		}
		return s.sendImage(ctx, target, text, *uploaded)
	default:
		uploaded, err := uploadMediaToCDNWithAPI(ctx, filePath, target.ToUserID, s.cdnBaseURL, UploadMediaTypeFile, s.api, s.httpClient, s.timeout)
		if err != nil {
			return "", err
		}
		return s.sendFile(ctx, target, text, filepath.Base(filePath), *uploaded)
	}
}

func buildTextMessageRequest(target Target, text, clientID string) SendMessageRequest {
	items := make([]MessageItem, 0, 1)
	if text != "" {
		items = append(items, MessageItem{
			Type:     MessageItemTypeText,
			TextItem: &TextItem{Text: text},
		})
	}
	return SendMessageRequest{
		Message: &WeixinMessage{
			ToUserID:     target.ToUserID,
			ClientID:     clientID,
			MessageType:  MessageTypeBot,
			MessageState: MessageStateFinish,
			ItemList:     items,
			ContextToken: target.ContextToken,
		},
	}
}

func buildItemMessageRequest(target Target, item MessageItem, clientID string) SendMessageRequest {
	return SendMessageRequest{
		Message: &WeixinMessage{
			ToUserID:     target.ToUserID,
			ClientID:     clientID,
			MessageType:  MessageTypeBot,
			MessageState: MessageStateFinish,
			ItemList:     []MessageItem{item},
			ContextToken: target.ContextToken,
		},
	}
}

func (s *sender) sendMediaItems(ctx context.Context, target Target, text string, mediaItem MessageItem) (string, error) {
	items := make([]MessageItem, 0, 2)
	if text != "" {
		items = append(items, MessageItem{Type: MessageItemTypeText, TextItem: &TextItem{Text: text}})
	}
	items = append(items, mediaItem)

	var lastID string
	for _, item := range items {
		lastID = GenerateID("openclaw-weixin")
		if err := s.api.SendMessage(ctx, buildItemMessageRequest(target, item, lastID), s.timeout); err != nil {
			return "", err
		}
	}
	return lastID, nil
}

func buildImageMessageItem(uploaded UploadedFileInfo) MessageItem {
	return MessageItem{
		Type: MessageItemTypeImage,
		ImageItem: &ImageItem{
			Media:   buildUploadedCDNMedia(uploaded),
			MidSize: uploaded.FileSizeCiphertext,
		},
	}
}

func buildVideoMessageItem(uploaded UploadedFileInfo) MessageItem {
	return MessageItem{
		Type: MessageItemTypeVideo,
		VideoItem: &VideoItem{
			Media:     buildUploadedCDNMedia(uploaded),
			VideoSize: uploaded.FileSizeCiphertext,
		},
	}
}

func buildFileMessageItem(fileName string, uploaded UploadedFileInfo) MessageItem {
	return MessageItem{
		Type: MessageItemTypeFile,
		FileItem: &FileItem{
			Media:    buildUploadedCDNMedia(uploaded),
			FileName: fileName,
			Length:   fmt.Sprintf("%d", uploaded.FileSize),
		},
	}
}

func buildUploadedCDNMedia(uploaded UploadedFileInfo) *CDNMedia {
	return &CDNMedia{
		EncryptQueryParam: uploaded.DownloadEncryptedQueryParam,
		AESKey:            base64.StdEncoding.EncodeToString([]byte(uploaded.AESKeyHex)),
		EncryptType:       1,
	}
}

func (t Target) validate(action string) error {
	if strings.TrimSpace(t.ToUserID) == "" {
		return fmt.Errorf("%s: toUserID is required", action)
	}
	if strings.TrimSpace(t.ContextToken) == "" {
		return fmt.Errorf("%s: contextToken is required", action)
	}
	return nil
}

func SetContextToken(accountID, userID, token string) {
	contextTokenStore.Lock()
	defer contextTokenStore.Unlock()
	contextTokenStore.values[accountID+":"+userID] = token
}

func GetContextToken(accountID, userID string) string {
	contextTokenStore.RLock()
	defer contextTokenStore.RUnlock()
	return contextTokenStore.values[accountID+":"+userID]
}

func IsMediaItem(item MessageItem) bool {
	return item.Type == MessageItemTypeImage ||
		item.Type == MessageItemTypeVideo ||
		item.Type == MessageItemTypeFile ||
		item.Type == MessageItemTypeVoice
}

func BodyFromItemList(items []MessageItem) string {
	for _, item := range items {
		if item.Type == MessageItemTypeText && item.TextItem != nil {
			text := item.TextItem.Text
			if item.RefMessage == nil {
				return text
			}
			if item.RefMessage.MessageItem != nil && IsMediaItem(*item.RefMessage.MessageItem) {
				return text
			}
			parts := make([]string, 0, 2)
			if item.RefMessage.Title != "" {
				parts = append(parts, item.RefMessage.Title)
			}
			if item.RefMessage.MessageItem != nil {
				if refBody := BodyFromItemList([]MessageItem{*item.RefMessage.MessageItem}); refBody != "" {
					parts = append(parts, refBody)
				}
			}
			if len(parts) == 0 {
				return text
			}
			return "[引用: " + strings.Join(parts, " | ") + "]\n" + text
		}
		if item.Type == MessageItemTypeVoice && item.VoiceItem != nil && item.VoiceItem.Text != "" {
			return item.VoiceItem.Text
		}
	}
	return ""
}

func WeixinMessageToContext(msg WeixinMessage, accountID string, opts *InboundMediaOptions) MessageContext {
	fromUserID := msg.FromUserID
	ctx := MessageContext{
		Body:          BodyFromItemList(msg.ItemList),
		From:          fromUserID,
		To:            fromUserID,
		AccountID:     accountID,
		OriginatingTo: fromUserID,
		MessageSID:    GenerateID("openclaw-weixin"),
		Timestamp:     msg.CreateTimeMS,
		Provider:      "openclaw-weixin",
		ChatType:      "direct",
		ContextToken:  msg.ContextToken,
	}
	if opts != nil {
		switch {
		case opts.DecryptedPicPath != "":
			ctx.MediaPath = opts.DecryptedPicPath
			ctx.MediaType = "image/*"
		case opts.DecryptedVideoPath != "":
			ctx.MediaPath = opts.DecryptedVideoPath
			ctx.MediaType = "video/mp4"
		case opts.DecryptedFilePath != "":
			ctx.MediaPath = opts.DecryptedFilePath
			ctx.MediaType = firstNonEmpty(opts.FileMediaType, "application/octet-stream")
		case opts.DecryptedVoicePath != "":
			ctx.MediaPath = opts.DecryptedVoicePath
			ctx.MediaType = firstNonEmpty(opts.VoiceMediaType, "audio/wav")
		}
	}
	return ctx
}

func GenerateID(prefix string) string {
	return fmt.Sprintf("%s:%d-%s", prefix, time.Now().UnixMilli(), randomHex(4))
}

func randomHex(n int) string {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("%x", time.Now().UnixNano())
	}
	return hex.EncodeToString(buf)
}
