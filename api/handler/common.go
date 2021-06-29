package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"

	"app/config"
	"app/provider/model"
	"app/response"
	"github.com/goapt/gee"
	"github.com/goapt/logger"
	"github.com/ilibs/identicon"
)

func getLoginUser(c *gee.Context) *model.Users {
	if u, ok := c.Get("userInfo"); ok {
		return u.(*model.Users)
	}
	return nil
}

type Common struct {
}

func NewCommon() *Common {
	return &Common{}
}

func (m *Common) Avatar(c *gee.Context) gee.Response {
	name := c.DefaultQuery("name", "default")

	// New Generator: Rehuse
	ig, err := identicon.New(
		"fifsky", // Namespace
		5,        // Number of blocks (Size)
		5,        // Density
	)

	if err != nil {
		panic(err) // Invalid Size or Density
	}

	ii, err := ig.Draw(name) // Generate an IdentIcon

	if err != nil {
		return nil
	}
	// Takes the size in pixels and any io.Writer
	_ = ii.Png(300, c.Writer) // 300px * 300px
	return nil
}

func (m *Common) DingMsg(c *gee.Context) gee.Response {

	tt := c.GetHeader("timestamp")
	sign := c.GetHeader("sign")
	algorithm := hmac.New(sha256.New, []byte(config.App.Common.DingAppSecret))
	algorithm.Write([]byte(tt + "\n" + config.App.Common.DingAppSecret))
	sign2 := base64.StdEncoding.EncodeToString(algorithm.Sum(nil))

	body, _ := ioutil.ReadAll(c.Request.Body)
	logger.Info("body:", string(body))

	if sign != sign2 {
		logger.Error("sign error", map[string]interface{}{
			"sign":  sign,
			"sign2": sign2,
			"tt":    tt,
		})
		return response.Fail(c, 202, "签名错误")
	}

	return c.JSON(map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content": "我还在开发哦……",
		},
	})
}
