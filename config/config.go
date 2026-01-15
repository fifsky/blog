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
}

type ossConf struct {
	AccessKey    string `yaml:"access_key"`
	AccessSecret string `yaml:"access_secret"`
	Endpoint     string `yaml:"endpoint"`
	Bucket       string `yaml:"bucket"`
}

type WeixinConf struct {
	Appid     string `yaml:"appid"`
	AppSecret string `yaml:"app_secret"`
	Token     string `yaml:"token"`
}

type Config struct {
	Env     string
	AppName string     `yaml:"app_name"`
	Common  common     `yaml:"common"`
	DB      DBConf     `yaml:"database"`
	OSS     ossConf    `yaml:"oss"`
	Weixin  WeixinConf `yaml:"weixin"`
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
