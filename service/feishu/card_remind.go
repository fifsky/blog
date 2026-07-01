package feishu

import (
	"context"
	"fmt"
	"strconv"
	"text/template"

	"app/config"
	"app/pkg/aesutil"
	apiv1 "app/proto/gen/api/v1"
	"app/store"

	"github.com/samber/lo"
)

// RemindMessage 提醒卡片消息，驱动卡片模板渲染
type RemindMessage struct {
	Content string // 提醒内容
	Time    string // 提醒时间
	Token   string // 回调按钮携带的 token
	Result  string // 操作结果文本，仅结果卡片使用
}

// RemindCard 提醒卡片处理器，合并卡片构建和回调处理
type RemindCard struct {
	tplBuilder
	remind apiv1.RemindServiceHTTPServer
	store  *store.Store
	conf   *config.Config
}

// NewRemindCard 创建提醒卡片处理器
func NewRemindCard(remind apiv1.RemindServiceHTTPServer, store *store.Store, conf *config.Config) *RemindCard {
	return &RemindCard{
		tplBuilder: tplBuilder{
			cardTpl:   template.Must(template.New("remind").Funcs(tplFuncs).Parse(remindCardTemplate)),
			resultTpl: template.Must(template.New("remindResult").Funcs(tplFuncs).Parse(remindResultCardTemplate)),
		},
		remind: remind,
		store:  store,
		conf:   conf,
	}
}

func (c *RemindCard) Actions() []string { return []string{"remind_completed", "remind_later"} }

// BuildCard 构建提醒卡片
func (c *RemindCard) BuildCard(msg RemindMessage) string { return c.execCard(msg) }

// BuildResultCard 构建提醒结果卡片
func (c *RemindCard) BuildResultCard(msg RemindMessage) string { return c.execResult(msg) }

// Handle 处理提醒卡片回调，返回结果卡片 JSON 和结果文本
func (c *RemindCard) Handle(ctx context.Context, action, token string) (string, string, error) {
	req := apiv1.RemindActionRequest_builder{Token: token}.Build()

	var result *apiv1.TextResponse
	var err error
	switch action {
	case "remind_completed":
		result, err = c.remind.Change(ctx, req)
	case "remind_later":
		result, err = c.remind.Delay(ctx, req)
	default:
		return "", "", nil
	}

	if err != nil {
		return "", "", fmt.Errorf("操作失败: %w", err)
	}

	id, err := aesutil.AesDecode(c.conf.Common.TokenSecret, token)
	if err != nil {
		return "", "", fmt.Errorf("token错误:%w", err)
	}

	remind, err := c.store.GetRemind(ctx, lo.Must(strconv.Atoi(id)))
	if err != nil {
		return "", "", fmt.Errorf("记录未找到:%w", err)
	}

	responseText := "操作完成"
	if result != nil {
		responseText = result.GetText()
	}

	msg := RemindMessage{
		Content: remind.Content,
		Time:    remind.NextTime.Format("2006-01-02 15:04"),
		Result:  responseText,
	}
	return c.BuildResultCard(msg), responseText, nil
}

// remindCardTemplate 提醒卡片模板，包含"标记完成"和"稍后提醒"两个回调按钮
const remindCardTemplate = `{
    "schema": "2.0",
    "config": {
        "update_multi": true
    },
    "body": {
        "direction": "vertical",
        "vertical_spacing": "12px",
        "padding": "20px 20px 20px 20px",
        "elements": [
            {
                "tag": "markdown",
                "content": "**今日待办事项**",
                "text_align": "left",
                "text_size": "heading",
                "margin": "0px 0px 0px 0px"
            },
            {
                "tag": "markdown",
                "content": {{.Time | json}},
                "text_align": "left",
                "text_size": "normal",
                "margin": "0px 0px 0px 0px",
                "icon": {
                    "tag": "standard_icon",
                    "token": "calendar_outlined",
                    "color": "grey"
                }
            },
            {
                "tag": "markdown",
                "content": {{.Content | json}},
                "text_align": "left",
                "text_size": "normal",
                "margin": "0px 0px 0px 0px"
            },
            {
                "tag": "column_set",
                "flex_mode": "stretch",
                "horizontal_spacing": "8px",
                "horizontal_align": "left",
                "columns": [
                    {
                        "tag": "column",
                        "width": "auto",
                        "elements": [
                            {
                                "tag": "button",
                                "text": {
                                    "tag": "plain_text",
                                    "content": "标记完成"
                                },
                                "type": "primary_filled",
                                "width": "fill",
                                "size": "medium",
                                "behaviors": [
                                    {
                                        "type": "callback",
                                        "value": {
                                            "action": "remind_completed",
                                            "token": {{.Token | json}}
                                        }
                                    }
                                ],
                                "margin": "4px 0px 4px 0px"
                            }
                        ],
                        "vertical_spacing": "8px",
                        "horizontal_align": "left",
                        "vertical_align": "top"
                    },
                    {
                        "tag": "column",
                        "width": "auto",
                        "elements": [
                            {
                                "tag": "button",
                                "text": {
                                    "tag": "plain_text",
                                    "content": "稍后提醒"
                                },
                                "type": "default",
                                "width": "fill",
                                "size": "medium",
                                "behaviors": [
                                    {
                                        "type": "callback",
                                        "value": {
                                            "action": "remind_later",
                                            "token": {{.Token | json}}
                                        }
                                    }
                                ],
                                "margin": "4px 0px 4px 0px"
                            }
                        ],
                        "vertical_spacing": "8px",
                        "horizontal_align": "left",
                        "vertical_align": "top"
                    }
                ],
                "margin": "0px 0px 0px 0px",
                "element_id": "action"
            }
        ]
    },
    "header": {
        "title": {
            "tag": "plain_text",
            "content": "个人事项提醒"
        },
        "subtitle": {
            "tag": "plain_text",
            "content": ""
        },
        "template": "wathet",
        "padding": "12px 12px 12px 12px"
    }
}`

// remindResultCardTemplate 提醒结果卡片模板，展示操作结果文本
const remindResultCardTemplate = `{
    "schema": "2.0",
    "config": {
        "update_multi": true
    },
    "body": {
        "direction": "vertical",
        "vertical_spacing": "12px",
        "padding": "20px 20px 20px 20px",
        "elements": [
            {
                "tag": "markdown",
                "content": "**今日待办事项**",
                "text_align": "left",
                "text_size": "heading",
                "margin": "0px 0px 0px 0px"
            },
            {
                "tag": "markdown",
                "content": {{.Time | json}},
                "text_align": "left",
                "text_size": "normal",
                "margin": "0px 0px 0px 0px",
                "icon": {
                    "tag": "standard_icon",
                    "token": "calendar_outlined",
                    "color": "grey"
                }
            },
            {
                "tag": "markdown",
                "content": {{.Content | json}},
                "text_align": "left",
                "text_size": "normal",
                "margin": "0px 0px 0px 0px"
            },
            {
                "tag": "div",
                "text": {
                    "tag": "plain_text",
                    "content": {{.Result | json}},
                    "text_size": "notation",
                    "text_align": "left",
                    "text_color": "default"
                },
                "icon": {
                    "tag": "standard_icon",
                    "token": "warning_outlined",
                    "color": "blue"
                },
                "margin": "0px 0px 0px 0px"
            }
        ]
    },
    "header": {
        "title": {
            "tag": "plain_text",
            "content": "个人事项提醒"
        },
        "subtitle": {
            "tag": "plain_text",
            "content": ""
        },
        "template": "wathet",
        "padding": "12px 12px 12px 12px"
    }
}`
