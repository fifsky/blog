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

func TestAdminCate_ListCreateUpdateDelete(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("cates")...)
		svc := NewCate(store.New(db))
		resp, err := svc.List(context.Background(), &emptypb.Empty{})
		if err != nil || len(resp.List) == 0 {
			t.Fatalf("unexpected err=%v list=%v", err, resp.List)
		}
		_, err = svc.Create(context.Background(), &adminv1.CateCreateRequest{Name: "demo", Domain: "demo", Desc: "demo"})
		if err != nil {
			t.Fatalf("unexpected err=%v", err)
		}
		_, err2 := svc.Update(context.Background(), &adminv1.CateUpdateRequest{Id: 1, Domain: "test2", Name: "test2", Desc: "test2"})
		if err2 != nil {
			t.Fatalf("unexpected err=%v", err2)
		}
		_, err3 := svc.Delete(context.Background(), &adminv1.CateDeleteRequest{Id: 3})
		if err3 != nil {
			t.Fatalf("unexpected err=%v", err3)
		}
	})
}
