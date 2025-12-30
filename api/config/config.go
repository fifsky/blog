package config

import (
	"app/pkg/jsonutil"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v3"
)

type common struct {
	RobotToken  string `yaml:"robot_token"`
	TokenSecret string `yaml:"token_secret"`
	NotifyUrl   string `yaml:"notify_url"`
	NotifyToken string `yaml:"notify_token"`
}

type ossConf struct {
	AccessKey    string `yaml:"access_key"`
	AccessSecret string `yaml:"access_secret"`
	Endpoint     string `yaml:"endpoint"`
	Bucket       string `yaml:"bucket"`
}

type Config struct {
	Env     string
	Path    string
	AppName string  `yaml:"app_name"`
	Common  common  `yaml:"common"`
	DB      DBConf  `yaml:"database"`
	OSS     ossConf `yaml:"oss"`
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
	fmt.Println(jsonutil.Encode(conf))

	return conf
}
