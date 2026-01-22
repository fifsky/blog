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
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("photos", "regions")...)
		svc := NewTravel(store.New(db))

		resp, err := svc.GetFootprints(context.Background(), &apiv1.GetFootprintsRequest{})
		if err != nil {
			t.Fatalf("GetFootprints failed: %v", err)
		}

		if len(resp.Provinces) == 0 {
			t.Fatal("Expected non-empty provinces list")
		}

		if len(resp.Cities) == 0 {
			t.Fatal("Expected non-empty cities list")
		}
	})
}

func TestTravel_ListCityPhotos(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("photos", "regions")...)
		svc := NewTravel(store.New(db))

		resp, err := svc.ListCityPhotos(context.Background(), &apiv1.ListCityPhotosRequest{RegionId: "310100"})
		if err != nil {
			t.Fatalf("ListCityPhotos failed: %v", err)
		}

		if len(resp.Photos) == 0 {
			t.Fatal("Expected non-empty photos list for Shanghai")
		}

		// Verify photo details
		photo := resp.Photos[0]
		if photo.Title == "" {
			t.Fatal("Expected photo title")
		}
		if photo.Src == "" {
			t.Fatal("Expected photo src")
		}
		if photo.Thumbnail == "" {
			t.Fatal("Expected photo thumbnail")
		}
	})
}
