package admin

import (
	"context"
	"testing"

	adminv1 "app/proto/gen/admin/v1"
	"app/store"
	"app/testutil"

	"github.com/goapt/dbunit"
)

func TestRegion_ListByParent(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("regions")...)
		svc := NewRegion(store.New(db))

		// Test list provinces (parent_id = 0)
		resp, err := svc.ListByParent(context.Background(), &adminv1.RegionListRequest{ParentId: 0})
		if err != nil {
			t.Fatalf("ListByParent failed: %v", err)
		}
		if len(resp.List) == 0 {
			t.Fatal("Expected non-empty province list")
		}

		// Test list cities (parent_id = province region_id)
		if len(resp.List) > 0 {
			cityResp, err := svc.ListByParent(context.Background(), &adminv1.RegionListRequest{ParentId: resp.List[0].RegionId})
			if err != nil {
				t.Fatalf("ListByParent for cities failed: %v", err)
			}
			// Zhejiang province should have cities
			if resp.List[0].RegionName == "浙江省" && len(cityResp.List) == 0 {
				t.Fatal("Expected non-empty city list for Zhejiang")
			}
		}
	})
}
