package admin

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"app/config"
	adminv1 "app/proto/gen/admin/v1"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var _ adminv1.OSSServiceServer = (*OSS)(nil)

type OSS struct {
	adminv1.UnimplementedOSSServiceServer
	conf *config.Config
}

func NewOSS(conf *config.Config) *OSS {
	return &OSS{conf: conf}
}

func (o *OSS) GetPresignURL(ctx context.Context, req *adminv1.GetPresignURLRequest) (*adminv1.GetPresignURLResponse, error) {
	// Extract region from endpoint (e.g., "oss-cn-shanghai.aliyuncs.com" -> "cn-shanghai")
	endpoint := o.conf.OSS.Endpoint
	region := ""
	if strings.HasPrefix(endpoint, "oss-") {
		parts := strings.Split(endpoint, ".")
		if len(parts) > 0 {
			region = strings.TrimPrefix(parts[0], "oss-")
		}
		region = strings.TrimSuffix(region, "-internal")
	}

	// Create OSS client using new SDK v2
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			o.conf.OSS.AccessKey,
			o.conf.OSS.AccessSecret,
		)).
		WithRegion(region).WithUseInternalEndpoint(false)

	client := oss.NewClient(cfg)

	// Generate upload path: blog/photos/${YYYY}/${MM}/${DD}/${filename}
	now := time.Now()
	ext := filepath.Ext(req.Filename)
	filename := fmt.Sprintf("%d%s", now.UnixNano(), ext)
	objectKey := fmt.Sprintf("blog/photos/%d/%02d/%02d/%s", now.Year(), now.Month(), now.Day(), filename)

	// Generate presigned PUT URL
	expiration := 10 * time.Minute
	result, err := client.Presign(ctx, &oss.PutObjectRequest{
		Bucket:      oss.Ptr(o.conf.OSS.Bucket),
		Key:         oss.Ptr(objectKey),
		ContentType: oss.Ptr("text/plain;charset=utf8"), // 请确保在服务端生成该签名URL时设置的ContentType与在使用URL时设置的ContentType一致
	}, oss.PresignExpires(expiration))
	if err != nil {
		return nil, fmt.Errorf("failed to generate presign URL: %w", err)
	}

	// CDN URL prefix
	cdnURL := fmt.Sprintf("https://static.fifsky.com/%s", objectKey)

	return &adminv1.GetPresignURLResponse{
		Url:    result.URL,
		CdnUrl: cdnURL,
	}, nil
}
