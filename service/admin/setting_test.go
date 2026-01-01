package admin

import (
	"context"
	"testing"

	adminv1 "app/proto/gen/admin/v1"
	"app/store"
	"app/testutil"

	"github.com/goapt/dbunit"
)

func TestAdminSetting_Update(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options")...)
		svc := NewSetting(store.New(db))
		resp, err := svc.Update(context.Background(), &adminv1.Options{Kv: map[string]string{"key": "abc", "key2": "def"}})
		if err != nil || len(resp.Kv) == 0 {
			t.Fatalf("unexpected err=%v kv=%v", err, resp.Kv)
		}
	})
}
