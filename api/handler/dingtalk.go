package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"strings"

	"app/config"
	"app/provider/model"
	"github.com/goapt/gee"
	"github.com/goapt/logger"
	"github.com/ilibs/gosql/v2"
	"github.com/tidwall/gjson"
)

// nolint
const (
	mooodTag  = "#心情#"
	remindTag = "#提醒#"
)

type DingTalk struct {
	db   *gosql.DB
	conf *config.Config
}

func NewDingTalk(db *gosql.DB, conf *config.Config) *DingTalk {
	return &DingTalk{db: db, conf: conf}
}

func (d *DingTalk) DingMsg(c *gee.Context) gee.Response {

	tt := c.GetHeader("timestamp")
	sign := c.GetHeader("sign")
	algorithm := hmac.New(sha256.New, []byte(d.conf.Common.DingAppSecret))
	algorithm.Write([]byte(tt + "\n" + d.conf.Common.DingAppSecret))
	sign2 := base64.StdEncoding.EncodeToString(algorithm.Sum(nil))

	body, _ := io.ReadAll(c.Request.Body)
	logger.Info("body:", string(body))

	if sign != sign2 {
		return d.dingReturn(c, "签名错误")
	}

	content := strings.TrimSpace(gjson.ParseBytes(body).Get("text.content").String())

	if content != "" {
		if strings.Contains(content, mooodTag) {
			mood := &model.Moods{
				UserId:  1,
				Content: strings.ReplaceAll(content, mooodTag, ""),
			}

			if _, err := d.db.Model(mood).Create(); err != nil {
				return d.dingReturn(c, "心情发表失败:"+err.Error())
			}

			return d.dingReturn(c, "心情发表成功")
		}
	}

	return d.dingReturn(c, "我还在开发哦……")
}

func (d *DingTalk) dingReturn(c *gee.Context, msg string) gee.Response {
	return c.JSON(gee.H{
		"msgtype": "text",
		"text": gee.H{
			"content": msg,
		},
	})
}
