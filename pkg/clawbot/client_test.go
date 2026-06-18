package clawbot

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoginInteractive(t *testing.T) {
	t.Parallel()

	var pollCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ilink/bot/get_bot_qrcode":
			_ = json.NewEncoder(w).Encode(map[string]string{
				"qrcode":             "qr-1",
				"qrcode_img_content": "weixin://qr-content",
			})
		case "/ilink/bot/get_qrcode_status":
			pollCount++
			if pollCount == 1 {
				_ = json.NewEncoder(w).Encode(map[string]string{
					"status": "scaned",
				})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]string{
				"status":        "confirmed",
				"bot_token":     "bot-token",
				"ilink_bot_id":  "demo@im.bot",
				"baseurl":       "https://returned.example",
				"ilink_user_id": "user@im.wechat",
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	tempDir := t.TempDir()
	client := NewClient(Options{
		BaseURL:           server.URL,
		HTTPClient:        server.Client(),
		PollInterval:      5 * time.Millisecond,
		QRLongPollTimeout: 100 * time.Millisecond,
	})

	account, err := client.LoginInteractive(context.Background(), InteractiveLoginOptions{
		Timeout: time.Second,
		SaveDir: tempDir,
	})
	if err != nil {
		t.Fatalf("LoginInteractive returned error: %v", err)
	}

	if account.AccountID != "demo@im.bot" {
		t.Fatalf("unexpected account id: %q", account.AccountID)
	}
	if account.BotToken != "bot-token" {
		t.Fatalf("unexpected bot token: %q", account.BotToken)
	}
	if account.UserID != "user@im.wechat" {
		t.Fatalf("unexpected user id: %q", account.UserID)
	}

	loaded, err := LoadAccount(tempDir, "demo@im.bot")
	if err != nil {
		t.Fatalf("LoadAccount returned error: %v", err)
	}
	if loaded.BotToken != "bot-token" {
		t.Fatalf("unexpected stored token: %q", loaded.BotToken)
	}
}

func TestWaitLoginRefreshesExpiredQR(t *testing.T) {
	t.Parallel()

	var qrFetches int
	var polls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ilink/bot/get_bot_qrcode":
			qrFetches++
			_ = json.NewEncoder(w).Encode(map[string]string{
				"qrcode":             "qr-refreshed",
				"qrcode_img_content": "weixin://qr-refreshed",
			})
		case "/ilink/bot/get_qrcode_status":
			polls++
			if polls == 1 {
				_ = json.NewEncoder(w).Encode(map[string]string{"status": "expired"})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]string{
				"status":       "confirmed",
				"bot_token":    "bot-token",
				"ilink_bot_id": "demo@im.bot",
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewClient(Options{
		BaseURL:           server.URL,
		HTTPClient:        server.Client(),
		PollInterval:      5 * time.Millisecond,
		QRLongPollTimeout: 100 * time.Millisecond,
	})

	session := &LoginSession{
		SessionKey: "session",
		QRCode:     "qr-old",
		QRContent:  "weixin://old",
		StartedAt:  time.Now(),
	}

	account, err := client.WaitLogin(context.Background(), session, WaitOptions{Timeout: time.Second})
	if err != nil {
		t.Fatalf("WaitLogin returned error: %v", err)
	}
	if account.AccountID != "demo@im.bot" {
		t.Fatalf("unexpected account id: %q", account.AccountID)
	}
	if qrFetches != 1 {
		t.Fatalf("expected one QR refresh, got %d", qrFetches)
	}
	if session.QRCode != "qr-refreshed" {
		t.Fatalf("expected session QR code to refresh, got %q", session.QRCode)
	}
}

func TestCheckLoginStatus(t *testing.T) {
	t.Parallel()

	t.Run("scaned returns current session", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/ilink/bot/get_qrcode_status" {
				t.Fatalf("unexpected path: %q", r.URL.Path)
			}
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "scaned"})
		}))
		defer server.Close()

		client := NewClient(Options{
			BaseURL:           server.URL,
			HTTPClient:        server.Client(),
			QRLongPollTimeout: 100 * time.Millisecond,
		})
		session := &LoginSession{
			SessionKey: "session",
			QRCode:     "qr-1",
			QRContent:  "weixin://qr-1",
			StartedAt:  time.Now(),
		}

		result, err := client.CheckLogin(context.Background(), session)
		if err != nil {
			t.Fatalf("CheckLogin returned error: %v", err)
		}
		if result.Status != "scaned" || result.Account != nil || result.Refreshed {
			t.Fatalf("unexpected CheckLogin result: %#v", result)
		}
	})

	t.Run("expired refreshes session for page qr", func(t *testing.T) {
		var qrFetches int
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/ilink/bot/get_qrcode_status":
				_ = json.NewEncoder(w).Encode(map[string]string{"status": "expired"})
			case "/ilink/bot/get_bot_qrcode":
				qrFetches++
				_ = json.NewEncoder(w).Encode(map[string]string{
					"qrcode":             "qr-2",
					"qrcode_img_content": "weixin://qr-2",
				})
			default:
				http.NotFound(w, r)
			}
		}))
		defer server.Close()

		client := NewClient(Options{
			BaseURL:           server.URL,
			HTTPClient:        server.Client(),
			QRLongPollTimeout: 100 * time.Millisecond,
		})
		session := &LoginSession{
			SessionKey: "session",
			QRCode:     "qr-1",
			QRContent:  "weixin://qr-1",
			StartedAt:  time.Now(),
		}

		result, err := client.CheckLogin(context.Background(), session)
		if err != nil {
			t.Fatalf("CheckLogin returned error: %v", err)
		}
		if result.Status != "wait" || !result.Refreshed || qrFetches != 1 {
			t.Fatalf("unexpected CheckLogin refresh result: %#v fetches=%d", result, qrFetches)
		}
		if session.QRCode != "qr-2" || session.QRContent != "weixin://qr-2" {
			t.Fatalf("expected refreshed session, got %#v", session)
		}
	})

	t.Run("confirmed returns account", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/ilink/bot/get_qrcode_status" {
				t.Fatalf("unexpected path: %q", r.URL.Path)
			}
			_ = json.NewEncoder(w).Encode(map[string]string{
				"status":        "confirmed",
				"bot_token":     "bot-token",
				"ilink_bot_id":  "demo@im.bot",
				"baseurl":       "https://returned.example",
				"ilink_user_id": "user@im.wechat",
			})
		}))
		defer server.Close()

		client := NewClient(Options{
			BaseURL:           server.URL,
			HTTPClient:        server.Client(),
			QRLongPollTimeout: 100 * time.Millisecond,
		})

		result, err := client.CheckLogin(context.Background(), &LoginSession{
			SessionKey: "session",
			QRCode:     "qr-1",
			QRContent:  "weixin://qr-1",
			StartedAt:  time.Now(),
		})
		if err != nil {
			t.Fatalf("CheckLogin returned error: %v", err)
		}
		if result.Status != "confirmed" || result.Account == nil {
			t.Fatalf("unexpected confirmed result: %#v", result)
		}
		if result.Account.AccountID != "demo@im.bot" || result.Account.BotToken != "bot-token" {
			t.Fatalf("unexpected account: %#v", result.Account)
		}
	})
}

func TestSaveAndListAccounts(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	if _, err := SaveAccount(tempDir, &Account{AccountID: "b@im.bot", BotToken: "b"}); err != nil {
		t.Fatalf("SaveAccount returned error: %v", err)
	}
	if _, err := SaveAccount(tempDir, &Account{AccountID: "a@im.bot", BotToken: "a"}); err != nil {
		t.Fatalf("SaveAccount returned error: %v", err)
	}

	accounts, err := ListAccounts(tempDir)
	if err != nil {
		t.Fatalf("ListAccounts returned error: %v", err)
	}
	if len(accounts) != 2 {
		t.Fatalf("expected 2 accounts, got %d", len(accounts))
	}
	if accounts[0].AccountID != "a@im.bot" || accounts[1].AccountID != "b@im.bot" {
		t.Fatalf("accounts not sorted: %#v", accounts)
	}
}

func TestSaveAccountUsesSafeFilename(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	filePath, err := SaveAccount(tempDir, &Account{AccountID: "demo@im.bot", BotToken: "token"})
	if err != nil {
		t.Fatalf("SaveAccount returned error: %v", err)
	}
	if filepath.Base(filePath) == "demo@im.bot.json" {
		t.Fatalf("expected safe encoded filename, got %q", filePath)
	}
	if _, err := os.Stat(filePath); err != nil {
		t.Fatalf("expected saved file: %v", err)
	}
}
