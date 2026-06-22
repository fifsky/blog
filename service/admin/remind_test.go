package admin

import (
	"context"
	"testing"

	"app/pkg/dbunit"
	adminv1 "app/proto/gen/admin/v1"
	"app/store"
	"app/testutil"
)

func TestAdminRemind_ListCreateDelete(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("reminds")...)
		svc := NewRemind(store.New(db))
		resp, err := svc.List(context.Background(), adminv1.RemindListRequest_builder{Page: 1}.Build())
		if err != nil || len(resp.GetList()) == 0 {
			t.Fatalf("unexpected err=%v list=%v", err, resp.GetList())
		}
		_, err = svc.Create(context.Background(), adminv1.RemindCreateRequest_builder{Cron: "* * * * *", Content: "demo"}.Build())
		if err != nil {
			t.Fatalf("unexpected err=%v", err)
		}
		_, err2 := svc.Delete(context.Background(), adminv1.RemindDeleteRequest_builder{Id: 8}.Build())
		if err2 != nil {
			t.Fatalf("unexpected err=%v", err2)
		}
	})
}
