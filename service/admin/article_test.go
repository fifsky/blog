package admin

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"app/config"
	"app/pkg/ossutil"
	adminv1 "app/proto/gen/admin/v1"
	"app/store"
	"app/store/model"
	"app/testutil"

	"github.com/goapt/dbunit"
	"go.uber.org/mock/gomock"
)

func TestAdminArticle_CreateDeleteUpload(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts")...)
		conf := &config.Config{}
		conf.OSS.Endpoint = "oss-cn-shanghai-internal.aliyuncs.com"
		conf.OSS.AccessKey = "test"
		conf.OSS.AccessSecret = "test"
		conf.OSS.Bucket = "test"
		svc := NewArticle(store.New(db), conf)

		// Create
		ctx := SetLoginUser(context.Background(), &model.User{Id: 1})
		resp, err := svc.Create(ctx, &adminv1.ArticleCreateRequest{CateId: 1, Type: 1, Title: "test", Url: "", Content: "test"})
		if err != nil || resp.Id == 0 {
			t.Fatalf("unexpected err=%v resp=%v", err, resp)
		}

		// Delete
		_, err = svc.Delete(context.Background(), &adminv1.ArticleDeleteRequest{Ids: []int32{4}})
		if err != nil {
			t.Fatalf("unexpected err=%v", err)
		}

		// Upload
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUploader := ossutil.NewMockUploader(ctrl)
		svc.upl = mockUploader
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
		respRR := httptest.NewRecorder()
		svc.Upload(respRR, req)
		if respRR.Code != http.StatusOK || !strings.Contains(respRR.Body.String(), ".png") {
			t.Fatalf("unexpected: code=%d body=%s", respRR.Code, respRR.Body.String())
		}
	})
}
