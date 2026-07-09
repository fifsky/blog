package openapi

import (
	"context"
	"testing"

	"app/pkg/dbunit"
	apiv1 "app/proto/gen/api/v1"
	"app/store"
	"app/testutil"
)

func TestMood_List(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("moods", "users"))
		svc := NewMood(store.New(db))
		resp, err := svc.List(context.Background(), apiv1.MoodListRequest_builder{Page: 1}.Build())
		if err != nil || len(resp.GetList()) == 0 {
			t.Fatalf("unexpected err=%v list=%v", err, resp.GetList())
		}
	})
}
