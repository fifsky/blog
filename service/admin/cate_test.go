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
		if err != nil || len(resp.GetList()) == 0 {
			t.Fatalf("unexpected err=%v list=%v", err, resp.GetList())
		}
		_, err = svc.Create(context.Background(), adminv1.CateCreateRequest_builder{Name: "demo", Domain: "demo", Desc: "demo"}.Build())
		if err != nil {
			t.Fatalf("unexpected err=%v", err)
		}
		_, err2 := svc.Update(context.Background(), adminv1.CateUpdateRequest_builder{Id: 1, Domain: "test2", Name: "test2", Desc: "test2"}.Build())
		if err2 != nil {
			t.Fatalf("unexpected err=%v", err2)
		}
		_, err3 := svc.Delete(context.Background(), adminv1.CateDeleteRequest_builder{Id: 3}.Build())
		if err3 != nil {
			t.Fatalf("unexpected err=%v", err3)
		}
	})
}
