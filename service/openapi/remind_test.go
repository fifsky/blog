package openapi

import (
	"context"
	"strconv"
	"testing"

	"app/config"
	"app/pkg/aesutil"
	"app/pkg/dbunit"
	apiv1 "app/proto/gen/api/v1"
	"app/store"
	"app/testutil"
)

func getRemindTestToken(id int, conf *config.Config) string {
	token, _ := aesutil.AesEncode(conf.Common.TokenSecret, strconv.Itoa(id))
	return token
}

func TestRemind_Change(t *testing.T) {
	conf := &config.Config{}
	conf.Common.TokenSecret = "abcdabcdabcdabcd"

	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("reminds")...)
		svc := NewRemind(store.New(db), conf)
		resp, err := svc.Change(context.Background(), apiv1.RemindActionRequest_builder{Token: getRemindTestToken(8, conf)}.Build())
		if err != nil || resp.GetText() == "" {
			t.Fatalf("unexpected err=%v resp=%v", err, resp)
		}
	})
}

func TestRemind_Delay(t *testing.T) {
	conf := &config.Config{}
	conf.Common.TokenSecret = "abcdabcdabcdabcd"

	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("reminds")...)
		svc := NewRemind(store.New(db), conf)
		resp, err := svc.Delay(context.Background(), apiv1.RemindActionRequest_builder{Token: getRemindTestToken(8, conf)}.Build())
		if err != nil || resp.GetText() == "" {
			t.Fatalf("unexpected err=%v resp=%v", err, resp)
		}
	})
}
