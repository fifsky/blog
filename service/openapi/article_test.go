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
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts")...)
		svc := NewArticle(store.New(db), nil)
		resp, err := svc.Archive(context.Background(), &emptypb.Empty{})
		if err != nil || len(resp.GetList()) == 0 {
			t.Fatalf("unexpected err=%v list=%v", err, resp.GetList())
		}
	})
}

func TestArticle_Calendar(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts")...)
		svc := NewArticle(store.New(db), nil)
		resp, err := svc.Calendar(context.Background(), apiv1.ArticleCalendarRequest_builder{Year: 2012, Month: 9}.Build())
		if err != nil {
			t.Fatalf("unexpected err=%v", err)
		}
		// 移除 type = 1 过滤后，返回所有已发布文章的日期
		if len(resp.GetDays()) != 2 {
			t.Fatalf("unexpected days=%v", resp.GetDays())
		}
	})
}

func TestArticle_List(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts", "users", "cates")...)
		svc := NewArticle(store.New(db), nil)
		resp, err := svc.List(context.Background(), apiv1.ArticleListRequest_builder{Page: 1}.Build())
		if err != nil || len(resp.GetList()) == 0 {
			t.Fatalf("unexpected err=%v list=%v", err, resp.GetList())
		}
	})
}

func TestArticle_List_Day(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts", "users", "cates")...)
		svc := NewArticle(store.New(db), nil)
		resp, err := svc.List(context.Background(), apiv1.ArticleListRequest_builder{Year: "2012", Month: "09", Day: "10", Page: 1}.Build())
		if err != nil || len(resp.GetList()) != 1 {
			t.Fatalf("unexpected err=%v list_len=%d", err, len(resp.GetList()))
		}
	})
}

func TestArticle_PrevNext(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts")...)
		svc := NewArticle(store.New(db), nil)
		resp, err := svc.PrevNext(context.Background(), apiv1.PrevNextRequest_builder{Id: 7}.Build())
		if err != nil || resp.GetPrev() == nil || resp.GetNext() == nil {
			t.Fatalf("unexpected err=%v prev=%v,next=%v", err, resp.GetPrev(), resp.GetNext())
		}
	})
}

func TestArticle_Detail(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts", "users", "cates")...)
		svc := NewArticle(store.New(db), nil)
		item, err := svc.Detail(context.Background(), apiv1.ArticleDetailRequest_builder{Id: 7}.Build())
		if err != nil || item == nil || item.GetId() == 0 {
			t.Fatalf("unexpected err=%v item=%v", err, item)
		}
	})
}

func TestArticle_Feed(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts", "users")...)
		svc := NewArticle(store.New(db), nil)
		resp, err := svc.Feed(context.Background(), &emptypb.Empty{})
		if err != nil || resp == nil || len(resp.Data) == 0 {
			t.Fatalf("unexpected err=%v resp=%v", err, resp)
		}
	})
}
