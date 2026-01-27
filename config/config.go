package config

import (
	"log"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v3"
)

type common struct {
	StoragePath string `yaml:"storage_path"`
	RobotToken  string `yaml:"robot_token"`
	TokenSecret string `yaml:"token_secret"`
	NotifyUrl   string `yaml:"notify_url"`
	NotifyToken string `yaml:"notify_token"`
	AIToken     string `yaml:"ai_token"`
	AIEndpoint  string `yaml:"ai_endpoint"`
	AIModel     string `yaml:"ai_model"`
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
	Appid      string `yaml:"appid"`
	AppSecret  string `yaml:"app_secret"`
	EncryptKey string `yaml:"encrypt_key"`
}

type Config struct {
	Env     string
	AppName string             `yaml:"app_name"`
	Common  common             `yaml:"common"`
	DB      DBConf             `yaml:"database"`
	OSS     ossConf            `yaml:"oss"`
	MCP     map[string]MCPConf `yaml:"mcp"`
	MiniAPP MiniAPPConf        `yaml:"miniapp"`
	Feishu  FeishuConf         `yaml:"feishu"`
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
