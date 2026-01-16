package openapi

import (
	"context"
	"testing"

	"app/store"
	"app/testutil"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/goapt/dbunit"
)

func TestSetting_Get(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options")...)
		svc := NewSetting(store.New(db))
		resp, err := svc.Get(context.Background(), &emptypb.Empty{})
		if err != nil || len(resp.Kv) == 0 {
			t.Fatalf("unexpected err=%v kv=%v", err, resp.Kv)
		}
	})
}
