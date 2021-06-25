package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goapt/acm"
	"github.com/goapt/envconf"
	"github.com/goapt/gee"
	"github.com/goapt/logger"
	"github.com/ilibs/gosql/v2"
)

type common struct {
	Debug         string `yml:"debug"`
	StoragePath   string `yml:"storage_path"`
	RobotToken    string `yml:"robot_token"`
	TokenSecret   string `yml:"token_secret"`
	DingAppSecret string `yml:"ding_app_secret"`
	EnvType       string `yml:"env_type"`
}

type ossConf struct {
	AccessKey    string `yml:"access_key"`
	AccessSecret string `yml:"access_secret"`
	Endpoint     string `yml:"endpoint"`
	Bucket       string `yml:"bucket"`
}

type acmConf struct {
	AccessKey    string `yml:"access_key"`
	AccessSecret string `yml:"access_secret"`
	Endpoint     string `yml:"endpoint"`
	RoleName     string `yml:"role_name"`
	Group        string `yml:"group"`
	Namespace    string `yml:"namespace"`
	DataId       string `yml:"data_id"`
}

type Config struct {
	Env       string
	Path      string
	AppName   string                   `yml:"app_name"`
	Common    common                   `yml:"common"`
	Log       logger.Config            `yml:"log"`
	DB        map[string]*gosql.Config `yml:"database"`
	OSS       ossConf                  `yml:"oss"`
	Acm       acmConf                  `yml:"acm"`
	StartTime time.Time
	IsTesting bool
}

var App = &Config{
	StartTime: time.Now(),
}

func init() {
	gee.ArgsInit()
	Load(gee.ExtCliArgs)
}

// 判断是否为测试执行
func isTestMode() bool {
	// test执行文件的路径后缀带.test，生产环境的可执行文件，不可能带.test后缀
	if strings.HasSuffix(os.Args[0], ".test") {
		return true
	}

	testVars := map[string]bool{
		"-test.v":   true,
		"-test.run": true,
	}

	for _, str := range os.Args {
		if testVars[str] {
			return true
		}
	}

	return false
}

func Load(args map[string]string) {
	appPath := args["config"]

	if appPath == "" {
		// 由于go test执行路径是临时目录，因此寻找配置文件要从编译路径查找
		if isTestMode() {
			App.IsTesting = true
			_, file, _, _ := runtime.Caller(0)
			appPath = filepath.Dir(filepath.Dir(file))
		} else {
			appPath = "./"
		}
	}

	App.Path = appPath

	conf, err := envconf.New(filepath.Join(appPath, "config.yml"))
	if err != nil {
		log.Fatalf("config error %s", err.Error())
	}

	if args["env"] != "" {
		conf.Set("env", args["env"])
	}

	if err := conf.Unmarshal(App); err != nil {
		log.Fatal("config unmarshal error:", err)
	}

	if App.Common.EnvType == "acm" && App.Env != "local" {
		acmconf := acm.NewAcm(func(c *acm.Acm) {
			c.SpasSecretKey = App.Acm.AccessKey
			c.SpasSecretKey = App.Acm.AccessSecret
			c.RoleName = App.Acm.RoleName
		})

		content, err := acmconf.Get(App.Acm.Namespace, App.Acm.Group, App.Acm.DataId)
		if err != nil {
			log.Fatal("acm config get error:", err)
		}
		// load env config
		if err := conf.EnvWithReader(strings.NewReader(content)); err != nil {
			log.Fatal("config env error:", err)
		}

	} else {
		if !App.IsTesting {
			// load config
			if err := conf.Env(filepath.Join(appPath, ".env")); err != nil {
				log.Fatal("config env error:", err)
			}
		}
	}

	if err := conf.Unmarshal(App); err != nil {
		log.Fatal("config unmarshal error:", err)
	}

	if !filepath.IsAbs(App.Common.StoragePath) {
		App.Common.StoragePath = filepath.Join(appPath, App.Common.StoragePath)
	}

	if App.Env == "local" && !App.IsTesting {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	// debug model
	if args["debug"] != "" {
		App.Common.Debug = args["debug"]
	}

	for _, d := range App.DB {
		d.ShowSql = args["show-sql"] == "on"
	}
}
