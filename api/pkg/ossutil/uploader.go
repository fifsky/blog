package ossutil

import (
	"context"
	"io"
	"net/http"

	"app/config"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

//go:generate mockgen --source=uploader.go -destination=mock_uploader.go -package=ossutil
type Uploader interface {
	Put(ctx context.Context, filename string, body io.Reader) error
}

type AliyunUploader struct {
	conf       *config.Config
	httpClient *http.Client
}

func NewAliyunUploader(conf *config.Config, httpClient *http.Client) *AliyunUploader {
	return &AliyunUploader{conf: conf, httpClient: httpClient}
}

func (u *AliyunUploader) Put(ctx context.Context, filename string, body io.Reader) error {
	client, err := oss.New(u.conf.OSS.Endpoint, u.conf.OSS.AccessKey, u.conf.OSS.AccessSecret, func(client *oss.Client) {
		client.HTTPClient = u.httpClient
	}, oss.AuthVersion(oss.AuthV4))
	if err != nil {
		return err
	}
	bucket, err := client.Bucket(u.conf.OSS.Bucket)
	if err != nil {
		return err
	}
	return bucket.PutObject(filename, body)
}
