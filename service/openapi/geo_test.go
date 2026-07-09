package openapi

import (
	"context"
	"testing"

	"app/pkg/dbunit"
	apiv1 "app/proto/gen/api/v1"
	"app/store"
	"app/testutil"
)

func TestGeo_GetNearestRegion(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("regions"))
		svc := NewGeo(store.New(db))

		// Test with valid coordinates (near Shanghai)
		resp, err := svc.GetNearestRegion(context.Background(), apiv1.GetNearestRegionRequest_builder{Latitude: 31.2304,
			Longitude: 121.4737}.Build(),
		)
		if err != nil {
			t.Fatalf("GetNearestRegion failed: %v", err)
		}

		if resp.GetCityId() == 0 {
			t.Fatal("Expected non-zero city ID")
		}

		if resp.GetCityName() == "" {
			t.Fatal("Expected non-empty city name")
		}

		if resp.GetProvinceId() == 0 {
			t.Fatal("Expected non-zero province ID")
		}

		if resp.GetProvinceName() == "" {
			t.Fatal("Expected non-empty province name")
		}
	})
}
