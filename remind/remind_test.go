package remind

import (
	"net/http"
	"os"
	"testing"

	"app/config"
	"app/model"
	"app/pkg/bark"
)

func Test_messageForBark(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	conf := config.Config{}
	conf.Common.NotifyUrl = os.Getenv("BARK_URL")
	conf.Common.NotifyToken = os.Getenv("BARK_TOKEN")
	conf.Common.TokenSecret = "2134ascd24"
	r := New(nil, &conf, nil, bark.New(http.DefaultClient, conf.Common.NotifyUrl, conf.Common.NotifyToken))
	r.messageForBark("这是一个此时这是一个此时这是一个此时这是一个此时这是一个此时这是一个此时这是一个此时这是一个此时", &model.Remind{Id: 1})
}
