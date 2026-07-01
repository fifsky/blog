package openapi

import (
	"context"
	"testing"

	"app/config"
	"app/pkg/dbunit"
	apiv1 "app/proto/gen/api/v1"
	"app/store"
	"app/testutil"
)

func TestUser_Login(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("users")...)
		conf := &config.Config{}
		conf.Common.TokenSecret = "abcdabcdabcdabcd"
		svc := NewUser(store.New(db), conf, nil, nil, nil)
		resp, err := svc.Login(context.Background(), apiv1.LoginRequest_builder{UserName: "test", Password: "test"}.Build())
		if err != nil || len(resp.GetAccessToken()) == 0 {
			t.Fatalf("unexpected err=%v token=%s", err, resp.GetAccessToken())
		}
		if resp.GetExpiresAt() == 0 {
			t.Fatalf("expected expires_at > 0, got %d", resp.GetExpiresAt())
		}
		_, err2 := svc.Login(context.Background(), apiv1.LoginRequest_builder{}.Build())
		if err2 == nil {
			t.Fatalf("expected validation error")
		}
		_, err3 := svc.Login(context.Background(), apiv1.LoginRequest_builder{UserName: "test", Password: "test234"}.Build())
		if err3 == nil {
			t.Fatalf("expected error for wrong password")
		}
		_, err4 := svc.Login(context.Background(), apiv1.LoginRequest_builder{UserName: "stop", Password: "test"}.Build())
		if err4 == nil {
			t.Fatalf("expected error for stopped user")
		}
	})
}
