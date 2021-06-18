package handler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"app/config"
	"app/model"
	"app/response"

	"github.com/goapt/gee"
	"github.com/ilibs/identicon"
	"github.com/tidwall/gjson"
)

func getLoginUser(c *gee.Context) *model.Users {
	if u, ok := c.Get("userInfo"); ok {
		return u.(*model.Users)
	}
	return nil
}

var Handle404 gee.HandlerFunc = func(c *gee.Context) gee.Response {
	return response.Fail(c, 404, "接口不存在")
}

var Avatar gee.HandlerFunc = func(c *gee.Context) gee.Response {
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
	ii.Png(300, c.Writer) // 300px * 300px
	return nil
}

func TCaptchaVerify(ticket, randstr, ip string) error {
	p := &url.Values{}
	p.Add("aid", config.App.Common.TCaptchaId)
	p.Add("AppSecretKey", config.App.Common.TCaptchaSecret)
	p.Add("Ticket", ticket)
	p.Add("Randstr", randstr)
	p.Add("UserIP", ip)

	fmt.Println("https://ssl.captcha.qq.com/ticket/verify?" + p.Encode())

	req, err := http.Get("https://ssl.captcha.qq.com/ticket/verify?" + p.Encode())

	if err != nil {
		return err
	}
	defer req.Body.Close()

	str, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	ret := gjson.ParseBytes(str)
	if ret.Get("response").Int() != 1 {
		return errors.New(ret.Get("err_msg").String())
	}

	return nil
}
