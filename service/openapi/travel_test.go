package openapi

import (
	"context"
	"testing"

	apiv1 "app/proto/gen/api/v1"
	"app/store"
	"app/testutil"

	"github.com/goapt/dbunit"
)

func TestTravel_GetFootprints(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("footprints")...)
		svc := NewTravel(store.New(db))

		resp, err := svc.GetFootprints(context.Background(), &apiv1.GetFootprintsRequest{})
		if err != nil {
			t.Fatalf("GetFootprints failed: %v", err)
		}

		if len(resp.Footprints) == 0 {
			t.Fatal("Expected non-empty footprints list")
		}

		fp := resp.Footprints[0]
		if fp.Name == "" {
			t.Fatal("Expected footprint name")
		}
		if fp.Longitude == "" || fp.Latitude == "" {
			t.Fatal("Expected footprint coordinates")
		}
	})
}
