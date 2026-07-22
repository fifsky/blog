package config

import (
	"app/pkg/agent"
	"app/pkg/doubaoasr"
	"app/pkg/litestream"
	"app/pkg/miniapp"
	"app/service/feishu"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"go.yaml.in/yaml/v3"
)

type common struct {
	StoragePath string `yaml:"storage_path"`
	TokenSecret string `yaml:"token_secret"`
	MCPToken    string `yaml:"mcp_token"`
}

type OssConfig struct {
	AccessKey    string `yaml:"access_key"`
	AccessSecret string `yaml:"access_secret"`
	Endpoint     string `yaml:"endpoint"`
	Bucket       string `yaml:"bucket"`
}

type Config struct {
	Env        string
	LogLevel   string                     `yaml:"log_level"`
	AppName    string                     `yaml:"app_name"`
	Common     common                     `yaml:"common"`
	DB         Database                   `yaml:"database"`
	OSS        OssConfig                  `yaml:"oss"`
	Litestream litestream.Config          `yaml:"litestream"`
	MCP        map[string]agent.MCPConfig `yaml:"mcp"`
	MiniAPP    miniapp.Config             `yaml:"miniapp"`
	Feishu     feishu.Config              `yaml:"feishu"`
	DoubaoASR  doubaoasr.Config           `yaml:"doubao_asr"`
}

func (c *Config) GetLogLevel() slog.Level {
	switch strings.ToLower(c.LogLevel) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func New() *Config {
	conf := &Config{}

	appPath := "./"
	content, err := os.ReadFile(filepath.Join(appPath, "config.yml"))
	if err != nil {
		log.Fatalf("config error %s", err.Error())
	}
	if err := yaml.Unmarshal(content, conf); err != nil {
		log.Fatal("config unmarshal error:", err)
	}

	_ = os.Setenv("APP_ENV", conf.Env)

	return conf
}
