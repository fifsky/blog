package handler

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"app/config"
	"app/pkg/ossutil"
	"app/store"
	"app/testutil"

	"github.com/goapt/dbunit"
	"go.uber.org/mock/gomock"
)

func TestArticle_Archive(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts")...)
		handler := NewArticle(store.New(db), nil)
		req := httptest.NewRequest(http.MethodPost, "/api/article/archive", bytes.NewReader([]byte(`{}`)))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler.Archive(rr, req)
		if rr.Code != http.StatusOK || !strings.Contains(rr.Body.String(), `"code":200`) {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
	})
}

func TestArticle_List(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts", "users", "cates")...)
		handler := NewArticle(store.New(db), nil)
		rr := doJSON(handler.List, "/api/article/list", map[string]any{"page": 1})
		if rr.Code != http.StatusOK || !strings.Contains(rr.Body.String(), `"list"`) {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
	})
}

func TestArticle_PrevNext(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts")...)
		handler := NewArticle(store.New(db), nil)
		rr := doJSON(handler.PrevNext, "/api/article/prevnext", map[string]any{"id": 7})
		if rr.Code != http.StatusOK || !strings.Contains(rr.Body.String(), `"prev"`) || !strings.Contains(rr.Body.String(), `"next"`) {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
	})
}

func TestArticle_Detail(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts", "users", "cates")...)
		handler := NewArticle(store.New(db), nil)
		rr := doJSON(handler.Detail, "/api/article/detail", map[string]any{"id": 7})
		if rr.Code != http.StatusOK || !strings.Contains(rr.Body.String(), `"code":200`) {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
	})
}

func TestArticle_Feed(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts", "users")...)
		handler := NewArticle(store.New(db), nil)
		rr := doJSON(handler.Feed, "/feed.xml", map[string]any{})
		if rr.Code != http.StatusOK || !strings.Contains(rr.Body.String(), "<feed") {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
	})
}

func TestArticle_Create(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts")...)
		handler := NewArticle(store.New(db), nil)
		// success with user context
		rr := doJSONWithUser(handler.Create, "/api/admin/article/create", map[string]any{"cate_id": 1, "type": 1, "title": "test", "url": "", "content": "test"})
		if rr.Code != http.StatusOK || !strings.Contains(rr.Body.String(), `"code":200`) {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
		// empty title error
		rr2 := doJSONWithUser(handler.Create, "/api/admin/article/create", map[string]any{"cate_id": 1, "type": 1, "title": "", "url": "", "content": "test"})
		if !strings.Contains(rr2.Body.String(), "title为必填字段") {
			t.Fatalf("unexpected body: %s", rr2.Body.String())
		}
	})
}

func TestArticle_Delete(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts")...)
		handler := NewArticle(store.New(db), nil)
		rr := doJSON(handler.Delete, "/api/admin/article/delete", map[string]any{"id": 4})
		if rr.Code != http.StatusOK || !strings.Contains(rr.Body.String(), `"code":200`) {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
	})
}

func TestArticle_Upload(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts")...)
		conf := &config.Config{}
		conf.OSS.Endpoint = "oss-cn-shanghai-internal.aliyuncs.com"
		conf.OSS.AccessKey = "test"
		conf.OSS.AccessSecret = "test"
		conf.OSS.Bucket = "test"

		handler := NewArticle(store.New(db), conf)
		handler.httpClient = http.DefaultClient
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUploader := ossutil.NewMockUploader(ctrl)
		handler.uploader = mockUploader
		mockUploader.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		rrf, err := os.Open(testutil.TestDataPath("go.png"))
		if err != nil {
			t.Fatalf("open file: %v", err)
		}
		defer rrf.Close()

		bb := &bytes.Buffer{}
		writer := multipart.NewWriter(bb)
		part, err := writer.CreateFormFile("uploadFile", "go.png")
		if err != nil {
			t.Fatalf("formfile: %v", err)
		}
		if _, err = io.Copy(part, rrf); err != nil {
			t.Fatalf("copy: %v", err)
		}
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/admin/article/upload", bb)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		resp := httptest.NewRecorder()
		handler.Upload(resp, req)
		if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), ".png") {
			t.Fatalf("unexpected: code=%d body=%s", resp.Code, resp.Body.String())
		}
	})
}
