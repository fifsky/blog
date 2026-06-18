package clawbot

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestClientUnifiedAPIHandlesBotOperations(t *testing.T) {
	t.Parallel()

	var paths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		switch r.URL.Path {
		case "/ilink/bot/getupdates":
			_ = json.NewEncoder(w).Encode(GetUpdatesResponse{
				Ret:           0,
				GetUpdatesBuf: "buf-next",
				Messages: []WeixinMessage{{
					FromUserID:   "user@im.wechat",
					ContextToken: "ctx-1",
					ItemList:     []MessageItem{{Type: MessageItemTypeText, TextItem: &TextItem{Text: "hello"}}},
				}},
			})
		case "/ilink/bot/sendtyping":
			var req SendTypingRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("decode typing: %v", err)
			}
			if req.ILinkUserID != "user@im.wechat" || req.TypingTicket != "ticket-1" || req.Status != TypingStatusTyping {
				t.Fatalf("unexpected typing request: %#v", req)
			}
			w.WriteHeader(http.StatusOK)
		case "/ilink/bot/sendmessage":
			var req SendMessageRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("decode message: %v", err)
			}
			if req.Message == nil || req.Message.ToUserID != "user@im.wechat" || req.Message.ContextToken != "ctx-1" {
				t.Fatalf("unexpected message request: %#v", req.Message)
			}
			w.WriteHeader(http.StatusOK)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewClient(Options{
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
	}).UseAccount(&Account{
		AccountID: "bot@im.bot",
		BotToken:  "bot-token",
		BaseURL:   server.URL,
	})

	updates, err := client.GetUpdates(context.Background(), GetUpdatesRequest{}, time.Second)
	if err != nil {
		t.Fatalf("GetUpdates returned error: %v", err)
	}
	if len(updates.Messages) != 1 || updates.GetUpdatesBuf != "buf-next" {
		t.Fatalf("unexpected updates: %#v", updates)
	}

	if err := client.SendTyping(context.Background(), SendTypingRequest{
		ILinkUserID:  "user@im.wechat",
		TypingTicket: "ticket-1",
		Status:       TypingStatusTyping,
	}, time.Second); err != nil {
		t.Fatalf("SendTyping returned error: %v", err)
	}

	if _, err := client.SendText(context.Background(), Target{
		ToUserID:     "user@im.wechat",
		ContextToken: "ctx-1",
	}, "hello back"); err != nil {
		t.Fatalf("SendText returned error: %v", err)
	}

	wantPaths := []string{"/ilink/bot/getupdates", "/ilink/bot/sendtyping", "/ilink/bot/sendmessage"}
	if !reflect.DeepEqual(paths, wantPaths) {
		t.Fatalf("paths = %#v, want %#v", paths, wantPaths)
	}
}

func TestClientListenUsesUnifiedClient(t *testing.T) {
	t.Parallel()

	var calls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		_ = json.NewEncoder(w).Encode(GetUpdatesResponse{
			Ret:           0,
			GetUpdatesBuf: "buf-next",
			Messages: []WeixinMessage{{
				FromUserID: "allowed",
				ItemList:   []MessageItem{{Type: MessageItemTypeText, TextItem: &TextItem{Text: "hello"}}},
			}},
		})
	}))
	defer server.Close()

	client := NewClient(Options{BaseURL: server.URL, HTTPClient: server.Client()})
	wantErr := context.Canceled
	ctx, cancel := context.WithCancel(context.Background())
	err := client.Listen(ctx, ListenOptions{
		AllowFrom: []string{"allowed"},
		OnMessages: func(ctx context.Context, messages []WeixinMessage) error {
			cancel()
			if len(messages) != 1 || messages[0].FromUserID != "allowed" {
				t.Fatalf("unexpected messages: %#v", messages)
			}
			return wantErr
		},
	})
	if err != wantErr {
		t.Fatalf("Listen error = %v, want %v", err, wantErr)
	}
	if calls != 1 {
		t.Fatalf("calls = %d, want 1", calls)
	}
}
