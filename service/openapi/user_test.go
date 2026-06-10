package openapi

import (
	"context"
	"testing"

	"app/config"
	apiv1 "app/proto/gen/api/v1"
	"app/store"
	"app/testutil"

	"github.com/goapt/dbunit"
)

func TestUser_Login(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("users")...)
		conf := &config.Config{}
		conf.Common.TokenSecret = "abcdabcdabcdabcd"
		svc := NewUser(store.New(db), conf)
		resp, err := svc.Login(context.Background(), apiv1.LoginRequest_builder{UserName: "test", Password: "test"}.Build())
		if err != nil || len(resp.GetAccessToken()) == 0 {
			t.Fatalf("unexpected err=%v token=%s", err, resp.GetAccessToken())
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
