package openapi

import (
	"context"
	"testing"

	"app/config"
	"app/pkg/dbunit"
	"app/store"
	"app/testutil"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestLink_All(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("links"))
		svc := NewLink(store.New(db), &config.Config{})
		resp, err := svc.All(context.Background(), &emptypb.Empty{})
		if err != nil || len(resp.GetList()) == 0 {
			t.Fatalf("unexpected err=%v list=%v", err, resp.GetList())
		}
	})
}
