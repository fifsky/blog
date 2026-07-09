package admin

import (
	"context"
	"testing"

	"app/pkg/dbunit"
	adminv1 "app/proto/gen/admin/v1"
	"app/store"
	"app/testutil"

	"google.golang.org/protobuf/types/known/emptypb"
)

func TestAdminLink_ListCreateDelete(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("links"))
		svc := NewLink(store.New(db))
		resp, err := svc.List(context.Background(), &emptypb.Empty{})
		if err != nil || len(resp.GetList()) == 0 {
			t.Fatalf("unexpected err=%v list=%v", err, resp.GetList())
		}
		_, err = svc.Create(context.Background(), adminv1.LinkCreateRequest_builder{Name: "demo", Url: "https://example.com", Desc: "demo"}.Build())
		if err != nil {
			t.Fatalf("unexpected err=%v", err)
		}
		_, err2 := svc.Delete(context.Background(), adminv1.LinkDeleteRequest_builder{Id: 3}.Build())
		if err2 != nil {
			t.Fatalf("unexpected err=%v", err2)
		}
	})
}
