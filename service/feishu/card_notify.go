package feishu

import "text/template"

// NotifyMessage 通用通知卡片消息
type NotifyMessage struct {
	Content string // 通知内容
	Time    string // 通知时间
}

// NotifyCard 通用通知卡片处理器，纯展示无回调操作
type NotifyCard struct {
	tplBuilder
}

// NewNotifyCard 创建通用通知卡片处理器
func NewNotifyCard() *NotifyCard {
	return &NotifyCard{
		tplBuilder: tplBuilder{
			cardTpl:   template.Must(template.New("notify").Funcs(tplFuncs).Parse(notifyCardTemplate)),
			resultTpl: nil, // 无回调，无结果卡片
		},
	}
}

// BuildCard 构建通用通知卡片
func (c *NotifyCard) BuildCard(msg NotifyMessage) string { return c.execCard(msg) }

// notifyCardTemplate 通用通知卡片模板，纯展示无按钮
const notifyCardTemplate = `{
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
                "content": {{.Content | json}},
                "text_align": "left",
                "text_size": "normal",
                "margin": "0px 0px 0px 0px"
            },
            {
                "tag": "markdown",
                "content": {{.Time | json}},
                "text_align": "left",
                "text_size": "notation",
                "margin": "0px 0px 0px 0px",
                "icon": {
                    "tag": "standard_icon",
                    "token": "calendar_outlined",
                    "color": "grey"
                }
            }
        ]
    },
    "header": {
        "title": {
            "tag": "plain_text",
            "content": "通知"
        },
        "subtitle": {
            "tag": "plain_text",
            "content": ""
        },
        "template": "indigo",
        "padding": "12px 12px 12px 12px"
    }
}`
