package openapi

import (
	"context"
	"strconv"
	"testing"

	"app/config"
	"app/pkg/aesutil"
	"app/pkg/wechat"
	apiv1 "app/proto/gen/api/v1"
	"app/store"
	"app/testutil"
	"github.com/goapt/dbunit"
)

func getRemindTestToken(id int, conf *config.Config) string {
	token, _ := aesutil.AesEncode(conf.Common.TokenSecret, strconv.Itoa(id))
	return token
}

func TestRemind_Change(t *testing.T) {
	conf := &config.Config{}
	conf.Common.TokenSecret = "abcdabcdabcdabcd"

	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.TestDSN, testutil.Schema(), testutil.Fixtures("reminds")...)
		robot := wechat.NewRobot("123")
		svc := NewRemind(store.New(db), wechat.NewRobot("123"), conf)
		resp, err := svc.Change(context.Background(), &apiv1.RemindActionRequest{Token: getRemindTestToken(8, conf)})
		if err != nil || resp.Text == "" {
			t.Fatalf("unexpected err=%v resp=%v", err, resp)
		}
		_ = robot
	})
}

func TestRemind_Delay(t *testing.T) {
	conf := &config.Config{}
	conf.Common.TokenSecret = "abcdabcdabcdabcd"

	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.TestDSN, testutil.Schema(), testutil.Fixtures("reminds")...)
		svc := NewRemind(store.New(db), wechat.NewRobot("123"), conf)
		resp, err := svc.Delay(context.Background(), &apiv1.RemindActionRequest{Token: getRemindTestToken(8, conf)})
		if err != nil || resp.Text == "" {
			t.Fatalf("unexpected err=%v resp=%v", err, resp)
		}
	})
}
