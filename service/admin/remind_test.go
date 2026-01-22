package admin

import (
	"context"
	"testing"

	adminv1 "app/proto/gen/admin/v1"
	"app/store"
	"app/testutil"

	"github.com/goapt/dbunit"
)

func TestAdminRemind_ListCreateDelete(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.TestDSN, testutil.Schema(), testutil.Fixtures("reminds")...)
		svc := NewRemind(store.New(db))
		resp, err := svc.List(context.Background(), &adminv1.RemindListRequest{Page: 1})
		if err != nil || len(resp.List) == 0 {
			t.Fatalf("unexpected err=%v list=%v", err, resp.List)
		}
		_, err = svc.Create(context.Background(), &adminv1.RemindCreateRequest{Type: 1, Content: "demo", Month: 1, Week: 0, Day: 1, Hour: 1, Minute: 1})
		if err != nil {
			t.Fatalf("unexpected err=%v", err)
		}
		_, err2 := svc.Delete(context.Background(), &adminv1.RemindDeleteRequest{Id: 8})
		if err2 != nil {
			t.Fatalf("unexpected err=%v", err2)
		}
	})
}
