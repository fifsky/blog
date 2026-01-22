package admin

import (
	"context"
	"testing"

	adminv1 "app/proto/gen/admin/v1"
	"app/store"
	"app/testutil"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/goapt/dbunit"
)

func TestAdminLink_ListCreateDelete(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.TestDSN, testutil.Schema(), testutil.Fixtures("links")...)
		svc := NewLink(store.New(db))
		resp, err := svc.List(context.Background(), &emptypb.Empty{})
		if err != nil || len(resp.List) == 0 {
			t.Fatalf("unexpected err=%v list=%v", err, resp.List)
		}
		_, err = svc.Create(context.Background(), &adminv1.LinkCreateRequest{Name: "demo", Url: "https://example.com", Desc: "demo"})
		if err != nil {
			t.Fatalf("unexpected err=%v", err)
		}
		_, err2 := svc.Delete(context.Background(), &adminv1.LinkDeleteRequest{Id: 3})
		if err2 != nil {
			t.Fatalf("unexpected err=%v", err2)
		}
	})
}
