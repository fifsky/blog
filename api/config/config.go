package config

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goapt/envconf"
	"github.com/goapt/gee"
	"github.com/goapt/logger"
	"github.com/ilibs/gosql/v2"
)

type common struct {
	Debug          string `yml:"debug"`
	StoragePath    string `yml:"storage_path"`
	RobotToken     string `yml:"robot_token"`
	SessionSecret  string `yml:"session_secret"`
	TokenSecret    string `yml:"token_secret"`
	TCaptchaId     string `yml:"tcaptcha_id"`
	TCaptchaSecret string `yml:"tcaptcha_secret"`
	DingAppSecret  string `yml:"ding_app_secret"`
}

type ossConf struct {
	AccessKey    string `yml:"access_key"`
	AccessSecret string `yml:"access_secret"`
	Endpoint     string `yml:"endpoint"`
	Bucket       string `yml:"bucket"`
}

type app struct {
	Env       string
	Path      string
	Common    common                   `yml:"common"`
	Log       logger.Config            `yml:"log"`
	DB        map[string]*gosql.Config `yml:"database"`
	OSS       ossConf                  `yml:"oss"`
	StartTime time.Time
	IsTesting bool
}

var App = &app{
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

	// load config
	if err := conf.Env(filepath.Join(appPath, ".env")); err != nil {
		log.Fatal("config env error:", err)
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

	// debug
	if App.Common.Debug == "on" {
		// log level
		App.Log.LogLevel = "debug"
		// log model
		App.Log.LogMode = "std"
	}

	if args["show-sql"] == "on" {
		for _, d := range App.DB {
			d.ShowSql = true
		}
	} else if args["show-sql"] == "off" {
		for _, d := range App.DB {
			d.ShowSql = false
		}
	}

	logger.Setting(func(c *logger.Config) {
		c.LogMode = App.Log.LogMode
		c.LogLevel = App.Log.LogLevel
		c.LogMaxFiles = App.Log.LogMaxFiles
		c.LogPath = filepath.Join(App.Common.StoragePath, "logs")
		c.LogSentryDSN = App.Log.LogSentryDSN
		c.LogSentryType = App.Log.LogSentryType
		c.LogDetail = App.Log.LogDetail
	})

	// test
	if os.Getenv("MYSQL_TEST_DSN") != "" {
		App.DB["default"].Dsn = os.Getenv("MYSQL_TEST_DSN")
	}
}

func ImportDB() ([]sql.Result, error) {
	sqlpath := "./blog.sql"
	rst, err := gosql.Import(sqlpath)
	return rst, err
}
