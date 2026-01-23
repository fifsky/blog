package ossutil

import (
	"context"
	"io"
	"strings"

	"app/config"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

//go:generate mockgen --source=uploader.go -destination=mock_uploader.go -package=ossutil
type Uploader interface {
	Put(ctx context.Context, filename string, body io.Reader) error
}

type AliyunUploader struct {
	conf   *config.Config
	client *oss.Client
}

func NewAliyunUploader(conf *config.Config) *AliyunUploader {
	// Extract region from endpoint (e.g., "oss-cn-shanghai.aliyuncs.com" -> "cn-shanghai")
	endpoint := conf.OSS.Endpoint
	region := ""
	if strings.HasPrefix(endpoint, "oss-") {
		parts := strings.Split(endpoint, ".")
		if len(parts) > 0 {
			region = strings.TrimPrefix(parts[0], "oss-")
		}
	}

	// Create OSS client using new SDK v2
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			conf.OSS.AccessKey,
			conf.OSS.AccessSecret,
		)).
		WithRegion(region)

	client := oss.NewClient(cfg)

	return &AliyunUploader{conf: conf, client: client}
}

func (u *AliyunUploader) Put(ctx context.Context, filename string, body io.Reader) error {
	// Create upload request
	putRequest := &oss.PutObjectRequest{
		Bucket: oss.Ptr(u.conf.OSS.Bucket),
		Key:    oss.Ptr(filename),
		Body:   body,
	}

	// Execute upload
	_, err := u.client.PutObject(ctx, putRequest)
	return err
}
