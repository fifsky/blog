package admin

import (
	"context"
	"testing"

	"app/pkg/dbunit"
	adminv1 "app/proto/gen/admin/v1"
	"app/store"
	"app/store/model"
	"app/testutil"
)

func TestAdminMood_CreateDelete(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("moods"))
		svc := NewMood(store.New(db))
		ctx := SetLoginUser(context.Background(), &model.User{Id: 1})
		_, err := svc.Create(ctx, adminv1.MoodCreateRequest_builder{Content: "demo"}.Build())
		if err != nil {
			t.Fatalf("unexpected err=%v", err)
		}
		_, err2 := svc.Delete(context.Background(), adminv1.MoodDeleteRequest_builder{Ids: []int32{1}}.Build())
		if err2 != nil {
			t.Fatalf("unexpected err=%v", err2)
		}
	})
}
