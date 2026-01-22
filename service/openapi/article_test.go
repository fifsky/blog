package openapi

import (
	"context"
	"testing"

	apiv1 "app/proto/gen/api/v1"
	"app/store"
	"app/testutil"

	"github.com/goapt/dbunit"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestArticle_Archive(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.TestDSN, testutil.Schema(), testutil.Fixtures("options", "posts")...)
		svc := NewArticle(store.New(db), nil)
		resp, err := svc.Archive(context.Background(), &emptypb.Empty{})
		if err != nil || len(resp.List) == 0 {
			t.Fatalf("unexpected err=%v list=%v", err, resp.List)
		}
	})
}

func TestArticle_Calendar(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.TestDSN, testutil.Schema(), testutil.Fixtures("options", "posts")...)
		svc := NewArticle(store.New(db), nil)
		resp, err := svc.Calendar(context.Background(), &apiv1.ArticleCalendarRequest{Year: 2012, Month: 9})
		if err != nil {
			t.Fatalf("unexpected err=%v", err)
		}
		if len(resp.Days) != 1 || resp.Days[0] != 10 {
			t.Fatalf("unexpected days=%v", resp.Days)
		}
	})
}

func TestArticle_List(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.TestDSN, testutil.Schema(), testutil.Fixtures("options", "posts", "users", "cates")...)
		svc := NewArticle(store.New(db), nil)
		resp, err := svc.List(context.Background(), &apiv1.ArticleListRequest{Page: 1})
		if err != nil || len(resp.List) == 0 {
			t.Fatalf("unexpected err=%v list=%v", err, resp.List)
		}
	})
}

func TestArticle_List_Day(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.TestDSN, testutil.Schema(), testutil.Fixtures("options", "posts", "users", "cates")...)
		svc := NewArticle(store.New(db), nil)
		resp, err := svc.List(context.Background(), &apiv1.ArticleListRequest{Year: "2012", Month: "09", Day: "10", Page: 1})
		if err != nil || len(resp.List) != 1 {
			t.Fatalf("unexpected err=%v list_len=%d", err, len(resp.List))
		}
	})
}

func TestArticle_PrevNext(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.TestDSN, testutil.Schema(), testutil.Fixtures("options", "posts")...)
		svc := NewArticle(store.New(db), nil)
		resp, err := svc.PrevNext(context.Background(), &apiv1.PrevNextRequest{Id: 7})
		if err != nil || resp.Prev == nil || resp.Next == nil {
			t.Fatalf("unexpected err=%v prev=%v,next=%v", err, resp.Prev, resp.Next)
		}
	})
}

func TestArticle_Detail(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.TestDSN, testutil.Schema(), testutil.Fixtures("options", "posts", "users", "cates")...)
		svc := NewArticle(store.New(db), nil)
		item, err := svc.Detail(context.Background(), &apiv1.ArticleDetailRequest{Id: 7})
		if err != nil || item == nil || item.Id == 0 {
			t.Fatalf("unexpected err=%v item=%v", err, item)
		}
	})
}

func TestArticle_Feed(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.TestDSN, testutil.Schema(), testutil.Fixtures("options", "posts", "users")...)
		svc := NewArticle(store.New(db), nil)
		resp, err := svc.Feed(context.Background(), &emptypb.Empty{})
		if err != nil || resp == nil || len(resp.Data) == 0 {
			t.Fatalf("unexpected err=%v resp=%v", err, resp)
		}
	})
}
