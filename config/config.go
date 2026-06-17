package config

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"go.yaml.in/yaml/v3"
)

type common struct {
	StoragePath string `yaml:"storage_path"`
	RobotToken  string `yaml:"robot_token"`
	TokenSecret string `yaml:"token_secret"`
	NotifyUrl   string `yaml:"notify_url"`
	NotifyToken string `yaml:"notify_token"`
	MCPToken    string `yaml:"mcp_token"`
}

type ossConf struct {
	AccessKey    string `yaml:"access_key"`
	AccessSecret string `yaml:"access_secret"`
	Endpoint     string `yaml:"endpoint"`
	Bucket       string `yaml:"bucket"`
}

type MCPConf struct {
	Name  string `yaml:"name"`  // MCP名称，用于前端展示
	URL   string `yaml:"url"`   // MCP StreamHttp地址
	Token string `yaml:"token"` // MCP Token
}

type MiniAPPConf struct {
	Appid     string `yaml:"appid"`
	AppSecret string `yaml:"app_secret"`
}

type FeishuConf struct {
	Appid                  string `yaml:"appid"`
	AppSecret              string `yaml:"app_secret"`
	UserID                 string `yaml:"user_id"`
	RemindTemplateID       string `yaml:"remind_template_id"`
	RemindResultTemplateID string `yaml:"remind_result_template_id"`
}

type Config struct {
	Env      string
	LogLevel string             `yaml:"log_level"`
	AppName  string             `yaml:"app_name"`
	Common   common             `yaml:"common"`
	DB       Database           `yaml:"database"`
	OSS      ossConf            `yaml:"oss"`
	MCP      map[string]MCPConf `yaml:"mcp"`
	MiniAPP  MiniAPPConf        `yaml:"miniapp"`
	Feishu   FeishuConf         `yaml:"feishu"`
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
