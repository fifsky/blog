package remind

import (
	"os"
	"testing"

	"app/config"
	"app/model"
)

func Test_messageForBark(t *testing.T) {
	// 跳过 Bark 消息推送测试，除非设置了相关的环境变量
	// 该测试依赖外部服务，通常只在本地开发或特定的集成测试环境中运行
	t.SkipNow()
	conf := config.Config{}
	conf.Common.NotifyUrl = os.Getenv("BARK_URL")
	conf.Common.NotifyToken = os.Getenv("BARK_TOKEN")
	conf.Common.TokenSecret = "2134ascd24"
	r := New(nil, &conf, nil)
	r.messageForBark("这是一个此时这是一个此时这是一个此时这是一个此时这是一个此时这是一个此时这是一个此时这是一个此时", &model.Remind{Id: 1})
}
