package config

import (
	"log"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v3"
)

type common struct {
	RobotToken  string `yml:"robot_token"`
	TokenSecret string `yml:"token_secret"`
	NotifyUrl   string `yml:"notify_url"`
	NotifyToken string `yml:"notify_token"`
}

type ossConf struct {
	AccessKey    string `yml:"access_key"`
	AccessSecret string `yml:"access_secret"`
	Endpoint     string `yml:"endpoint"`
	Bucket       string `yml:"bucket"`
}

type Config struct {
	Env     string
	Path    string
	AppName string  `yml:"app_name"`
	Common  common  `yml:"common"`
	DB      DBConf  `yml:"database"`
	OSS     ossConf `yml:"oss"`
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
	return conf
}
