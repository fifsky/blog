package admin

import (
	"context"
	"testing"

	adminv1 "app/proto/gen/admin/v1"
	"app/store"
	"app/store/model"
	"app/testutil"

	"github.com/goapt/dbunit"
)

func TestAdminMood_CreateDelete(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixture("moods"))
		svc := NewMood(store.New(db))
		ctx := SetLoginUser(context.Background(), &model.User{Id: 1})
		_, err := svc.Create(ctx, &adminv1.MoodCreateRequest{Content: "demo"})
		if err != nil {
			t.Fatalf("unexpected err=%v", err)
		}
		_, err2 := svc.Delete(context.Background(), &adminv1.MoodDeleteRequest{Ids: []int32{1}})
		if err2 != nil {
			t.Fatalf("unexpected err=%v", err2)
		}
	})
}
