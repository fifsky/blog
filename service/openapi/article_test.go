package openapi

import (
	"context"
	"encoding/xml"
	"testing"

	"app/pkg/dbunit"
	apiv1 "app/proto/gen/api/v1"
	"app/store"
	"app/testutil"

	"google.golang.org/protobuf/types/known/emptypb"
)

func TestArticle_Archive(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts"))
		svc := NewArticle(store.New(db))
		resp, err := svc.Archive(context.Background(), &emptypb.Empty{})
		if err != nil || len(resp.GetList()) == 0 {
			t.Fatalf("unexpected err=%v list=%v", err, resp.GetList())
		}
	})
}

func TestArticle_Calendar(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts"))
		svc := NewArticle(store.New(db))
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
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts", "users", "cates"))
		svc := NewArticle(store.New(db))
		resp, err := svc.List(context.Background(), apiv1.ArticleListRequest_builder{Page: 1}.Build())
		if err != nil || len(resp.GetList()) == 0 {
			t.Fatalf("unexpected err=%v list=%v", err, resp.GetList())
		}
	})
}

func TestArticle_List_Day(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts", "users", "cates"))
		svc := NewArticle(store.New(db))
		resp, err := svc.List(context.Background(), apiv1.ArticleListRequest_builder{Year: "2012", Month: "09", Day: "10", Page: 1}.Build())
		if err != nil || len(resp.GetList()) != 1 {
			t.Fatalf("unexpected err=%v list_len=%d", err, len(resp.GetList()))
		}
	})
}

func TestArticle_PrevNext(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts"))
		svc := NewArticle(store.New(db))
		resp, err := svc.PrevNext(context.Background(), apiv1.PrevNextRequest_builder{Id: 7}.Build())
		if err != nil || resp.GetPrev() == nil || resp.GetNext() == nil {
			t.Fatalf("unexpected err=%v prev=%v,next=%v", err, resp.GetPrev(), resp.GetNext())
		}
	})
}

func TestArticle_Detail(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts", "users", "cates"))
		svc := NewArticle(store.New(db))
		item, err := svc.Detail(context.Background(), apiv1.ArticleDetailRequest_builder{Id: 7}.Build())
		if err != nil || item == nil || item.GetId() == 0 {
			t.Fatalf("unexpected err=%v item=%v", err, item)
		}
	})
}

func TestArticle_Feed(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts", "users"))
		svc := NewArticle(store.New(db))
		resp, err := svc.Feed(context.Background(), &emptypb.Empty{})
		if err != nil || resp == nil || len(resp.Data) == 0 {
			t.Fatalf("unexpected err=%v resp=%v", err, resp)
		}
	})
}

func TestArticle_FeedUsesPostPublishedTimestamps(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts", "users"))
		svc := NewArticle(store.New(db))
		resp, err := svc.Feed(context.Background(), &emptypb.Empty{})
		if err != nil {
			t.Fatalf("unexpected err=%v", err)
		}
		if resp.GetContentType() != "application/atom+xml; charset=utf-8" {
			t.Fatalf("unexpected content type=%q", resp.GetContentType())
		}

		var feed atomFeed
		if err := xml.Unmarshal(resp.GetData(), &feed); err != nil {
			t.Fatalf("unmarshal atom feed: %v", err)
		}
		if feed.Updated != "2012-10-28T23:30:39+08:00" {
			t.Fatalf("unexpected feed updated=%q", feed.Updated)
		}
		if len(feed.Links) == 0 || feed.Links[0].Rel != "self" || feed.Links[0].Href != "https://api.fifsky.com/blog/feed.xml" {
			t.Fatalf("unexpected feed links=%v", feed.Links)
		}

		entries := make(map[string]atomEntry, len(feed.Entries))
		for _, entry := range feed.Entries {
			entries[entry.Title] = entry
		}

		example := entries["example"]
		if example.Updated != "2012-10-28T23:30:39+08:00" {
			t.Fatalf("unexpected example updated=%q", example.Updated)
		}
		if example.ID != "https://fifsky.com/article/8" {
			t.Fatalf("unexpected example id=%q", example.ID)
		}

		about := entries["关于"]
		if about.Updated != "2012-09-28T23:30:39+08:00" {
			t.Fatalf("unexpected about updated=%q", about.Updated)
		}
		if about.ID != "https://fifsky.com/article/7" {
			t.Fatalf("unexpected about id=%q", about.ID)
		}
	})
}

func TestArticle_FeedWithoutPostsHasUpdated(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options"))
		svc := NewArticle(store.New(db))
		resp, err := svc.Feed(context.Background(), &emptypb.Empty{})
		if err != nil {
			t.Fatalf("unexpected err=%v", err)
		}

		var feed atomFeed
		if err := xml.Unmarshal(resp.GetData(), &feed); err != nil {
			t.Fatalf("unmarshal atom feed: %v", err)
		}
		if feed.Updated == "" {
			t.Fatal("feed updated should not be empty")
		}
	})
}

type atomFeed struct {
	Updated string      `xml:"updated"`
	Links   []atomLink  `xml:"link"`
	Entries []atomEntry `xml:"entry"`
}

type atomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
}

type atomEntry struct {
	Title   string `xml:"title"`
	Updated string `xml:"updated"`
	ID      string `xml:"id"`
}
