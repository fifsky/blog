package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"

	"app/response"
	"github.com/goapt/gee"
	"github.com/goapt/logger"

	"app/config"
)

var DingMsg gee.HandlerFunc = func(c *gee.Context) gee.Response {

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
