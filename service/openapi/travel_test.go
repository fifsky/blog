package openapi

import (
	"context"
	"testing"

	"app/pkg/dbunit"
	apiv1 "app/proto/gen/api/v1"
	"app/store"
	"app/testutil"
)

func TestTravel_GetFootprints(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("footprints")...)
		svc := NewTravel(store.New(db))

		resp, err := svc.GetFootprints(context.Background(), apiv1.GetFootprintsRequest_builder{}.Build())
		if err != nil {
			t.Fatalf("GetFootprints failed: %v", err)
		}

		if len(resp.GetFootprints()) == 0 {
			t.Fatal("Expected non-empty footprints list")
		}

		fp := resp.GetFootprints()[0]
		if fp.GetName() == "" {
			t.Fatal("Expected footprint name")
		}
		if fp.GetLongitude() == "" || fp.GetLatitude() == "" {
			t.Fatal("Expected footprint coordinates")
		}
	})
}
