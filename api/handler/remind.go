package handler

import (
	"time"

	"github.com/goapt/gee"
	"github.com/goapt/golib/robot"
	"github.com/ilibs/gosql/v2"

	"app/model"
)

var RemindChange gee.HandlerFunc = func(c *gee.Context) gee.Response {
	remind, ok := c.Get("remind")

	if !ok {
		return c.Fail(202, "记录未找到")
	}

	v := remind.(*model.Reminds)

	_, err := gosql.Model(&model.Reminds{Status: 1}).Where("id = ?", v.Id).Update()

	if err != nil {
		return c.Fail(203, err)
	}

	_ = robot.Message("已确认收到提醒")
	return c.String("已确认收到提醒")
}

var RemindDelay gee.HandlerFunc = func(c *gee.Context) gee.Response {
	remind, ok := c.Get("remind")

	if !ok {
		return c.Fail(202, "记录未找到")
	}

	v := remind.(*model.Reminds)

	nextTime := time.Now().Add(10 * time.Minute)
	_, err := gosql.Model(&model.Reminds{NextTime: nextTime}).Where("id = ?", v.Id).Update()

	if err != nil {
		return c.Fail(203, err)
	}
	_ = robot.Message("将在10分钟后再次提醒")
	return c.String("将在10分钟后再次提醒")
}
