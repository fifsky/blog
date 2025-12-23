package config

import (
	"log"
	"path/filepath"
	"time"

	"github.com/goapt/envconf"
	"github.com/goapt/gee"
	"github.com/goapt/logger"
	"github.com/ilibs/gosql/v2"
)

type common struct {
	StoragePath   string `yml:"storage_path"`
	RobotToken    string `yml:"robot_token"`
	TokenSecret   string `yml:"token_secret"`
	DingAppSecret string `yml:"ding_app_secret"`
	NotifyUrl     string `yml:"notify_url"`
	NotifyToken   string `yml:"notify_token"`
}

type ossConf struct {
	AccessKey    string `yml:"access_key"`
	AccessSecret string `yml:"access_secret"`
	Endpoint     string `yml:"endpoint"`
	Bucket       string `yml:"bucket"`
}

type tencentCloud struct {
	SecretId  string `yml:"secret_id"`
	SecretKey string `yml:"secret_key"`
}

type Config struct {
	Env          string
	Path         string
	AppName      string                   `yml:"app_name"`
	Common       common                   `yml:"common"`
	Log          logger.Config            `yml:"log"`
	DB           map[string]*gosql.Config `yml:"database"`
	OSS          ossConf                  `yml:"oss"`
	TencentCloud tencentCloud             `yml:"tencent_cloud"`
	StartTime    time.Time
}

func New() *Config {
	gee.ArgsInit()

	conf := &Config{
		StartTime: time.Now(),
	}

	conf.load(gee.ExtCliArgs)
	return conf
}

func (c *Config) load(args map[string]string) {
	appPath := args["config"]
	if appPath == "" {
		appPath = "./"
	}

	c.Path = appPath

	conf, err := envconf.New(filepath.Join(appPath, "config.yml"))
	if err != nil {
		log.Fatalf("config error %s", err.Error())
	}

	// load config
	if err := conf.Env(filepath.Join(appPath, ".env")); err != nil {
		log.Fatal("config env error:", err)
	}

	if err := conf.Unmarshal(c); err != nil {
		log.Fatal("config unmarshal error:", err)
	}

	if !filepath.IsAbs(c.Common.StoragePath) {
		c.Common.StoragePath = filepath.Join(appPath, c.Common.StoragePath)
	}

	if c.Env == "local" {
		gee.SetMode(gee.DebugMode)
	} else {
		gee.SetMode(gee.ReleaseMode)
	}

	for _, d := range c.DB {
		d.ShowSql = args["show-sql"] == "on"
	}
}
