package service

import (
	"context"
	"testing"

	"app/store"
	"app/testutil"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/goapt/dbunit"
)

func TestLink_All(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("links")...)
		svc := NewLink(store.New(db))
		resp, err := svc.All(context.Background(), &emptypb.Empty{})
		if err != nil || len(resp.List) == 0 {
			t.Fatalf("unexpected err=%v list=%v", err, resp.List)
		}
	})
}
