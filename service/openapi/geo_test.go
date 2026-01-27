package openapi

import (
	"context"
	"testing"

	apiv1 "app/proto/gen/api/v1"
	"app/store"
	"app/testutil"

	"github.com/goapt/dbunit"
)

func TestGeo_GetNearestRegion(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("regions")...)
		svc := NewGeo(store.New(db))

		// Test with valid coordinates (near Shanghai)
		resp, err := svc.GetNearestRegion(context.Background(), &apiv1.GetNearestRegionRequest{
			Latitude:  31.2304,
			Longitude: 121.4737,
		})
		if err != nil {
			t.Fatalf("GetNearestRegion failed: %v", err)
		}

		if resp.CityId == 0 {
			t.Fatal("Expected non-zero city ID")
		}

		if resp.CityName == "" {
			t.Fatal("Expected non-empty city name")
		}

		if resp.ProvinceId == 0 {
			t.Fatal("Expected non-zero province ID")
		}

		if resp.ProvinceName == "" {
			t.Fatal("Expected non-empty province name")
		}
	})
}
