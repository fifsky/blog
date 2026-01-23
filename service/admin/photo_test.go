package admin

import (
	"context"
	"testing"

	adminv1 "app/proto/gen/admin/v1"
	"app/store"
	"app/testutil"

	"github.com/goapt/dbunit"
)

func TestPhoto_ListCreateUpdateDelete(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("photos", "regions")...)
		svc := NewPhoto(store.New(db))

		// Test List
		resp, err := svc.List(context.Background(), &adminv1.PhotoListRequest{Page: 1})
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if len(resp.List) == 0 {
			t.Fatal("Expected non-empty list")
		}

		// Test Create
		createResp, err := svc.Create(context.Background(), &adminv1.PhotoCreateRequest{
			Title:       "测试标题",
			Description: "测试描述",
			Srcs:        []string{"https://static.fifsky.com/blog/photos/2026/01/22/test.jpg"},
			Province:    "310000",
			City:        "310100",
		})
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		if createResp.Id == 0 {
			t.Fatal("Expected non-zero ID")
		}

		// Test Update
		_, err = svc.Update(context.Background(), &adminv1.PhotoUpdateRequest{
			Id:    createResp.Id,
			Title: "更新后的标题",
		})
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		// Test Delete
		_, err = svc.Delete(context.Background(), &adminv1.PhotoDeleteRequest{Id: createResp.Id})
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}
	})
}
