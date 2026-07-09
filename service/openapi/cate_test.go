package openapi

import (
	"context"
	"testing"

	"app/pkg/dbunit"
	"app/store"
	"app/testutil"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestCate_All(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("cates"))
		svc := NewCate(store.New(db))
		resp, err := svc.All(context.Background(), &emptypb.Empty{})
		if err != nil || len(resp.GetList()) == 0 {
			t.Fatalf("unexpected err=%v list=%v", err, resp.GetList())
		}
	})
}
