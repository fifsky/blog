package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"app/config"
	adminv1 "app/proto/gen/admin/v1"
	"app/store"
	"app/testutil"

	"github.com/goapt/dbunit"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestClawBotLoginPersistsAccount(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options")...)
		s := store.New(db)
		botToken := strings.Repeat("token-", 50)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/ilink/bot/get_bot_qrcode":
				require.Equal(t, "3", r.URL.Query().Get("bot_type"))
				_ = json.NewEncoder(w).Encode(map[string]string{
					"qrcode":             "qr-1",
					"qrcode_img_content": "weixin://qr-1",
				})
			case "/ilink/bot/get_qrcode_status":
				_ = json.NewEncoder(w).Encode(map[string]string{
					"status":        "confirmed",
					"bot_token":     botToken,
					"ilink_bot_id":  "demo@im.bot",
					"baseurl":       "https://returned.example",
					"ilink_user_id": "user@im.wechat",
				})
			default:
				http.NotFound(w, r)
			}
		}))
		defer server.Close()

		svc := NewClawBot(s, &config.Config{}, nil,
			WithClawBotBaseURL(server.URL),
			WithClawBotHTTPClient(server.Client()),
			WithClawBotMonitor(false),
		)

		session, err := svc.StartLogin(context.Background(), adminv1.ClawBotStartLoginRequest_builder{}.Build())
		require.NoError(t, err)
		require.Equal(t, "wait", session.GetStatus())
		require.Equal(t, "weixin://qr-1", session.GetQrContent())

		resp, err := svc.CheckLogin(context.Background(), adminv1.ClawBotCheckLoginRequest_builder{SessionKey: session.GetSessionKey()}.Build())
		require.NoError(t, err)
		require.True(t, resp.GetConnected())
		require.Equal(t, "confirmed", resp.GetStatus())
		require.Equal(t, "demo@im.bot", resp.GetAccount().GetAccountId())

		opts, err := s.GetOptions(context.Background())
		require.NoError(t, err)
		require.Equal(t, "demo@im.bot", opts["clawbot_account_id"])
		require.NotEmpty(t, opts["clawbot_bot_token_0"])
		require.Equal(t, "https://returned.example", opts["clawbot_base_url"])
		require.Equal(t, "user@im.wechat", opts["clawbot_user_id"])

		account, err := svc.loadAccount(context.Background())
		require.NoError(t, err)
		require.Equal(t, botToken, account.BotToken)

		status, err := svc.Status(context.Background(), &emptypb.Empty{})
		require.NoError(t, err)
		require.True(t, status.GetConnected())
		require.False(t, status.GetMonitoring())
		require.Equal(t, "demo@im.bot", status.GetAccount().GetAccountId())
	})
}

func TestClawBotCheckLoginRefreshesQRCode(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options")...)
		s := store.New(db)
		pollCount := 0

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/ilink/bot/get_bot_qrcode":
				if pollCount == 0 {
					_ = json.NewEncoder(w).Encode(map[string]string{
						"qrcode":             "qr-1",
						"qrcode_img_content": "weixin://qr-1",
					})
					return
				}
				_ = json.NewEncoder(w).Encode(map[string]string{
					"qrcode":             "qr-2",
					"qrcode_img_content": "weixin://qr-2",
				})
			case "/ilink/bot/get_qrcode_status":
				pollCount++
				_ = json.NewEncoder(w).Encode(map[string]string{"status": "expired"})
			default:
				http.NotFound(w, r)
			}
		}))
		defer server.Close()

		svc := NewClawBot(s, &config.Config{}, nil,
			WithClawBotBaseURL(server.URL),
			WithClawBotHTTPClient(server.Client()),
			WithClawBotMonitor(false),
		)

		session, err := svc.StartLogin(context.Background(), adminv1.ClawBotStartLoginRequest_builder{}.Build())
		require.NoError(t, err)

		resp, err := svc.CheckLogin(context.Background(), adminv1.ClawBotCheckLoginRequest_builder{SessionKey: session.GetSessionKey()}.Build())
		require.NoError(t, err)
		require.False(t, resp.GetConnected())
		require.Equal(t, "wait", resp.GetStatus())
		require.Equal(t, "weixin://qr-2", resp.GetSession().GetQrContent())
	})
}

func TestClawBotDisconnectClearsOptions(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options")...)
		s := store.New(db)
		_, err := s.UpdateOptions(context.Background(), map[string]string{
			"clawbot_account_id": "demo@im.bot",
			"clawbot_bot_token":  "bot-token",
			"clawbot_base_url":   "https://returned.example",
			"clawbot_user_id":    "user@im.wechat",
			"clawbot_saved_at":   "2026-06-18T00:00:00Z",
		})
		require.NoError(t, err)

		svc := NewClawBot(s, &config.Config{}, nil, WithClawBotMonitor(false))
		status, err := svc.Disconnect(context.Background(), &emptypb.Empty{})
		require.NoError(t, err)
		require.False(t, status.GetConnected())

		opts, err := s.GetOptions(context.Background())
		require.NoError(t, err)
		require.Empty(t, opts["clawbot_account_id"])
		require.Empty(t, opts["clawbot_bot_token"])
	})
}
