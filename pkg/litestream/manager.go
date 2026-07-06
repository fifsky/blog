// Package litestream 将 Litestream 作为 Go library 嵌入应用，实现 SQLite 实时流式备份到阿里云 OSS。
// 参考：https://litestream.io/guides/go-library/
package litestream

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"app/config"

	"github.com/benbjohnson/litestream"
	"github.com/benbjohnson/litestream/s3"
)

// Manager 管理 Litestream 复制生命周期。
// 用法：
//  1. 创建 Manager
//  2. 调用 Start(ctx) 初始化（内部会从 OSS 自动恢复数据库，若本地无文件）
//  3. 此后才打开应用的 sql.DB 连接
//  4. 应用退出前：先关闭 sql.DB，再调用 Stop() 关闭 Store
type Manager struct {
	store      *litestream.Store
	db         *litestream.DB
	replicaURI string // 用于日志展示的备份地址
}

// New 创建 Litestream 管理器
// 备份路径自动根据环境区分：开发环境 blog/sqlite/local，线上 blog/sqlite/prod
func New(conf *config.Config, dbPath string) *Manager {
	rPath := replicaPath(conf)

	// 创建 S3 客户端，通过自定义 Endpoint 对接阿里云 OSS（虚拟主机式访问）
	client := s3.NewReplicaClient()
	client.Bucket = conf.Litestream.Bucket
	client.Path = rPath
	client.Endpoint = conf.Litestream.Endpoint
	client.Region = conf.Litestream.Region
	client.AccessKeyID = conf.OSS.AccessKey
	client.SecretAccessKey = conf.OSS.AccessSecret
	// 阿里云 OSS 要求虚拟主机式访问（bucket.oss-cn-xxx.aliyuncs.com），
	// 不能使用 ForcePathStyle，保持默认 false

	// 包装 SQLite 数据库文件（注意：此时尚未打开，仅记录路径）
	db := litestream.NewDB(dbPath)

	// 创建副本并关联到数据库
	replica := litestream.NewReplicaWithClient(db, client)
	db.Replica = replica

	// 配置压缩级别（使用 Litestream 默认推荐配置）
	levels := litestream.DefaultCompactionLevels

	store := litestream.NewStore([]*litestream.DB{db}, levels)

	return &Manager{
		store:      store,
		db:         db,
		replicaURI: fmt.Sprintf("s3://%s/%s", conf.Litestream.Bucket, rPath),
	}
}

// Start 初始化复制：先尝试从 OSS 恢复（若本地无数据库文件），再启动后台复制。
// 应在打开应用的 sql.DB 之前调用。
func (m *Manager) Start(ctx context.Context) error {
	if m.store == nil {
		return fmt.Errorf("[litestream] store not initialized")
	}

	// 若本地数据库文件不存在，自动从 OSS 备份恢复
	if err := m.db.EnsureExists(ctx); err != nil {
		log.Printf("[litestream] restore skipped: %s", err)
	}

	// 打开 Store 启动后台复制协程
	if err := m.store.Open(ctx); err != nil {
		return fmt.Errorf("[litestream] failed to open store: %w", err)
	}

	log.Printf("[litestream] replication started → %s", m.replicaURI)
	return nil
}

// Stop 优雅关闭复制：先同步一次确保最新数据已上传，再关闭 Store。
// 调用前务必先关闭应用的 sql.DB。
func (m *Manager) Stop() error {
	if m.store == nil {
		return nil
	}

	// 优雅关闭：等待正在进行中的复制完成（最长 30 秒）
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := m.store.Close(ctx); err != nil {
		log.Printf("[litestream] store close error: %s", err)
		return err
	}

	log.Printf("[litestream] replication stopped")
	return nil
}

// replicaPath 根据环境生成备份路径
// 开发环境（非 prod）: blog/sqlite/local
// 线上环境（prod）:    blog/sqlite/prod
func replicaPath(conf *config.Config) string {
	env := strings.ToLower(conf.Env)
	suffix := "local"
	if env == "prod" || env == "production" {
		suffix = "prod"
	}
	return strings.TrimRight(conf.Litestream.Path, "/") + "/" + suffix
}
