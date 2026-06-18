package clawbot

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	DefaultBaseURL           = "https://ilinkai.weixin.qq.com"
	DefaultCDNBaseURL        = "https://novac2c.cdn.weixin.qq.com/c2c"
	DefaultBotType           = "3"
	DefaultChannelVersion    = "go-port"
	DefaultQRSessionTTL      = 5 * time.Minute
	DefaultQRLongPollTimeout = 35 * time.Second
	DefaultLoginTimeout      = 8 * time.Minute
	DefaultPollInterval      = time.Second
	DefaultMaxQRRefresh      = 3
)

type Options struct {
	BaseURL           string
	CDNBaseURL        string
	BotType           string
	Token             string
	AccountID         string
	Account           *Account
	RouteTag          string
	ChannelVersion    string
	HTTPClient        *http.Client
	Timeout           time.Duration
	QRSessionTTL      time.Duration
	QRLongPollTimeout time.Duration
	PollInterval      time.Duration
	MaxQRRefresh      int
}

type Client struct {
	baseURL           string
	cdnBaseURL        string
	botType           string
	token             string
	accountID         string
	routeTag          string
	channelVersion    string
	httpClient        *http.Client
	timeout           time.Duration
	qrSessionTTL      time.Duration
	qrLongPollTimeout time.Duration
	pollInterval      time.Duration
	maxQRRefresh      int
	configManager     *configManager
}

type LoginSession struct {
	SessionKey   string
	AccountHint  string
	QRCode       string
	QRContent    string
	StartedAt    time.Time
	RefreshCount int
}

type InteractiveLoginOptions struct {
	AccountHint string
	Timeout     time.Duration
	SaveDir     string
}

type WaitOptions struct {
	Timeout time.Duration
	SaveDir string
}

type Account struct {
	AccountID string `json:"account_id"`
	BotToken  string `json:"bot_token"`
	BaseURL   string `json:"base_url,omitempty"`
	UserID    string `json:"user_id,omitempty"`
	SavedAt   string `json:"saved_at,omitempty"`
}

type LoginCheckResult struct {
	Status    string
	Account   *Account
	Refreshed bool
}

type qrCodeResponse struct {
	QRCode       string `json:"qrcode"`
	QRCodeImgRaw string `json:"qrcode_img_content"`
}

type qrStatusResponse struct {
	Status    string `json:"status"`
	BotToken  string `json:"bot_token"`
	AccountID string `json:"ilink_bot_id"`
	BaseURL   string `json:"baseurl"`
	UserID    string `json:"ilink_user_id"`
}

func NewClient(opts Options) *Client {
	baseURL := strings.TrimSpace(opts.BaseURL)
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	cdnBaseURL := strings.TrimSpace(opts.CDNBaseURL)
	if cdnBaseURL == "" {
		cdnBaseURL = DefaultCDNBaseURL
	}

	botType := strings.TrimSpace(opts.BotType)
	if botType == "" {
		botType = DefaultBotType
	}

	channelVersion := strings.TrimSpace(opts.ChannelVersion)
	if channelVersion == "" {
		channelVersion = DefaultChannelVersion
	}

	httpClient := opts.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = defaultAPITimeout
	}

	qrSessionTTL := opts.QRSessionTTL
	if qrSessionTTL <= 0 {
		qrSessionTTL = DefaultQRSessionTTL
	}

	qrLongPollTimeout := opts.QRLongPollTimeout
	if qrLongPollTimeout <= 0 {
		qrLongPollTimeout = DefaultQRLongPollTimeout
	}

	pollInterval := opts.PollInterval
	if pollInterval <= 0 {
		pollInterval = DefaultPollInterval
	}

	maxQRRefresh := opts.MaxQRRefresh
	if maxQRRefresh <= 0 {
		maxQRRefresh = DefaultMaxQRRefresh
	}

	client := &Client{
		baseURL:           baseURL,
		cdnBaseURL:        cdnBaseURL,
		botType:           botType,
		token:             strings.TrimSpace(opts.Token),
		accountID:         strings.TrimSpace(opts.AccountID),
		routeTag:          strings.TrimSpace(opts.RouteTag),
		channelVersion:    channelVersion,
		httpClient:        httpClient,
		timeout:           timeout,
		qrSessionTTL:      qrSessionTTL,
		qrLongPollTimeout: qrLongPollTimeout,
		pollInterval:      pollInterval,
		maxQRRefresh:      maxQRRefresh,
	}
	if opts.Account != nil {
		client.UseAccount(opts.Account)
	}
	return client
}

func (c *Client) UseAccount(account *Account) *Client {
	if c == nil || account == nil {
		return c
	}
	c.accountID = strings.TrimSpace(account.AccountID)
	c.token = strings.TrimSpace(account.BotToken)
	if baseURL := strings.TrimSpace(account.BaseURL); baseURL != "" {
		c.baseURL = baseURL
	}
	c.configManager = nil
	return c
}

func (c *Client) StartLogin(ctx context.Context, accountHint string) (*LoginSession, error) {
	resp, err := c.fetchQRCode(ctx)
	if err != nil {
		return nil, err
	}

	return &LoginSession{
		SessionKey:   randomSessionKey(),
		AccountHint:  strings.TrimSpace(accountHint),
		QRCode:       resp.QRCode,
		QRContent:    resp.QRCodeImgRaw,
		StartedAt:    time.Now(),
		RefreshCount: 1,
	}, nil
}

func (c *Client) LoginInteractive(ctx context.Context, opts InteractiveLoginOptions) (*Account, error) {
	session, err := c.StartLogin(ctx, opts.AccountHint)
	if err != nil {
		return nil, err
	}

	return c.WaitLogin(ctx, session, WaitOptions{
		Timeout: opts.Timeout,
		SaveDir: opts.SaveDir,
	})
}

func (c *Client) CheckLogin(ctx context.Context, session *LoginSession) (*LoginCheckResult, error) {
	if session == nil {
		return nil, fmt.Errorf("login session is nil")
	}
	if time.Since(session.StartedAt) > c.qrSessionTTL {
		return nil, fmt.Errorf("login session expired")
	}

	status, err := c.pollQRStatus(ctx, session.QRCode)
	if err != nil {
		return nil, err
	}

	switch status.Status {
	case "wait", "scaned":
		return &LoginCheckResult{Status: status.Status}, nil
	case "expired":
		if err := c.refreshLoginSession(ctx, session); err != nil {
			return nil, err
		}
		return &LoginCheckResult{Status: "wait", Refreshed: true}, nil
	case "confirmed":
		account, err := c.accountFromStatus(status)
		if err != nil {
			return nil, err
		}
		return &LoginCheckResult{Status: "confirmed", Account: account}, nil
	default:
		return nil, fmt.Errorf("unexpected QR status %q", status.Status)
	}
}

func (c *Client) WaitLogin(ctx context.Context, session *LoginSession, opts WaitOptions) (*Account, error) {
	if session == nil {
		return nil, fmt.Errorf("login session is nil")
	}
	if time.Since(session.StartedAt) > c.qrSessionTTL {
		return nil, fmt.Errorf("login session expired")
	}

	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = DefaultLoginTimeout
	}
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		status, err := c.pollQRStatus(ctx, session.QRCode)
		if err != nil {
			return nil, err
		}

		switch status.Status {
		case "wait":
		case "scaned":
		case "expired":
			if err := c.refreshLoginSession(ctx, session); err != nil {
				return nil, err
			}
		case "confirmed":
			account, err := c.accountFromStatus(status)
			if err != nil {
				return nil, err
			}
			if opts.SaveDir != "" {
				if _, err := SaveAccount(opts.SaveDir, account); err != nil {
					return nil, fmt.Errorf("save account: %w", err)
				}
			}
			return account, nil
		default:
			return nil, fmt.Errorf("unexpected QR status %q", status.Status)
		}

		if err := sleepContext(ctx, c.pollInterval); err != nil {
			return nil, err
		}
	}

	return nil, fmt.Errorf("login timeout after %s", timeout)
}

func (c *Client) fetchQRCode(ctx context.Context) (*qrCodeResponse, error) {
	endpoint, err := joinURL(c.baseURL, "/ilink/bot/get_bot_qrcode")
	if err != nil {
		return nil, err
	}

	query := endpoint.Query()
	query.Set("bot_type", c.botType)
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.routeTag != "" {
		req.Header.Set("SKRouteTag", c.routeTag)
	}

	var resp qrCodeResponse
	if err := c.doJSON(req, &resp); err != nil {
		return nil, fmt.Errorf("fetch QR code: %w", err)
	}
	if resp.QRCode == "" || resp.QRCodeImgRaw == "" {
		return nil, fmt.Errorf("fetch QR code: empty QR payload")
	}
	return &resp, nil
}

func (c *Client) refreshLoginSession(ctx context.Context, session *LoginSession) error {
	if session.RefreshCount <= 0 {
		session.RefreshCount = 1
	}
	session.RefreshCount++
	if session.RefreshCount > c.maxQRRefresh {
		return fmt.Errorf("login timeout: QR code expired too many times")
	}
	refreshed, err := c.fetchQRCode(ctx)
	if err != nil {
		return fmt.Errorf("refresh QR code: %w", err)
	}
	session.QRCode = refreshed.QRCode
	session.QRContent = refreshed.QRCodeImgRaw
	session.StartedAt = time.Now()
	return nil
}

func (c *Client) accountFromStatus(status *qrStatusResponse) (*Account, error) {
	if status.AccountID == "" {
		return nil, fmt.Errorf("login confirmed but ilink_bot_id is missing")
	}
	return &Account{
		AccountID: status.AccountID,
		BotToken:  status.BotToken,
		BaseURL:   firstNonEmpty(status.BaseURL, c.baseURL),
		UserID:    status.UserID,
		SavedAt:   time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func (c *Client) pollQRStatus(ctx context.Context, qrCode string) (*qrStatusResponse, error) {
	endpoint, err := joinURL(c.baseURL, "/ilink/bot/get_qrcode_status")
	if err != nil {
		return nil, err
	}

	query := endpoint.Query()
	query.Set("qrcode", qrCode)
	endpoint.RawQuery = query.Encode()

	pollCtx, cancel := context.WithTimeout(ctx, c.qrLongPollTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(pollCtx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("iLink-App-ClientVersion", "1")
	if c.routeTag != "" {
		req.Header.Set("SKRouteTag", c.routeTag)
	}

	var resp qrStatusResponse
	if err := c.doJSON(req, &resp); err != nil {
		if errorsIsTimeout(err) {
			return &qrStatusResponse{Status: "wait"}, nil
		}
		return nil, fmt.Errorf("poll QR status: %w", err)
	}
	return &resp, nil
}

func (c *Client) doJSON(req *http.Request, out any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode JSON: %w", err)
	}
	return nil
}

func joinURL(base, pathPart string) (*url.URL, error) {
	baseURL, err := url.Parse(strings.TrimRight(base, "/") + "/")
	if err != nil {
		return nil, fmt.Errorf("parse base URL: %w", err)
	}
	ref, err := url.Parse(strings.TrimLeft(pathPart, "/"))
	if err != nil {
		return nil, fmt.Errorf("parse path: %w", err)
	}
	return baseURL.ResolveReference(ref), nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func randomSessionKey() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("session-%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(buf)
}

func sleepContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func errorsIsTimeout(err error) bool {
	if err == nil {
		return false
	}
	if strings.Contains(err.Error(), "context deadline exceeded") {
		return true
	}
	return false
}
