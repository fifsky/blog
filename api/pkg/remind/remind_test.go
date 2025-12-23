package remind

import (
	"os"
	"testing"

	"app/config"
	"app/provider/model"
)

func Test_messageForBark(t *testing.T) {
	t.SkipNow()
	conf := config.Config{}
	conf.Common.NotifyUrl = os.Getenv("BARK_URL")
	conf.Common.NotifyToken = os.Getenv("BARK_TOKEN")
	conf.Common.TokenSecret = "2134ascd24"
	messageForBark("这是一个此时这是一个此时这是一个此时这是一个此时这是一个此时这是一个此时这是一个此时这是一个此时", &model.Reminds{Id: 1}, &conf)
}
