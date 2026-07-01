package feishu

import (
	"context"
	"fmt"
	"strconv"
	"text/template"

	"app/config"
	"app/pkg/aesutil"
	"app/store"
	"app/store/model"
)

// LinkActionValue 友情链接卡片回调的 actionValue 结构体
type LinkActionValue struct {
	Action string `json:"action"` // 回调动作：link_approve / link_reject
	Token  string `json:"token"`  // 加密的业务 ID
}

// LinkMessage 友情链接卡片消息，驱动卡片模板渲染
type LinkMessage struct {
	Name   string // 站点名称
	URL    string // 站点地址
	Desc   string // 站点描述
	Token  string // 回调按钮携带的 token
	Result string // 操作结果文本，仅结果卡片使用
}

// LinkCard 友情链接卡片处理器，合并卡片构建和回调处理
type LinkCard struct {
	tplBuilder
	store *store.Store
	conf  *config.Config
}

// NewLinkCard 创建友情链接卡片处理器
func NewLinkCard(store *store.Store, conf *config.Config) *LinkCard {
	return &LinkCard{
		tplBuilder: tplBuilder{
			cardTpl:   template.Must(template.New("link").Funcs(tplFuncs).Parse(linkCardTemplate)),
			resultTpl: template.Must(template.New("linkResult").Funcs(tplFuncs).Parse(linkResultCardTemplate)),
		},
		store: store,
		conf:  conf,
	}
}

func (c *LinkCard) Actions() []string { return []string{"link_approve", "link_reject"} }

// BuildCard 构建友情链接审核卡片
func (c *LinkCard) BuildCard(msg LinkMessage) string { return c.execCard(msg) }

// BuildResultCard 构建友情链接审核结果卡片
func (c *LinkCard) BuildResultCard(msg LinkMessage) string { return c.execResult(msg) }

// Handle 处理友情链接卡片回调，返回结果卡片 JSON 和结果文本
func (c *LinkCard) Handle(ctx context.Context, actionValue map[string]any) (string, string, error) {
	actionVal, err := parseActionValue[LinkActionValue](actionValue)
	if err != nil {
		return "", "", fmt.Errorf("解析actionValue失败: %w", err)
	}

	id, err := aesutil.AesDecode(c.conf.Common.TokenSecret, actionVal.Token)
	if err != nil {
		return "", "", fmt.Errorf("token错误:%w", err)
	}

	linkID, err := strconv.Atoi(id)
	if err != nil {
		return "", "", fmt.Errorf("链接ID错误:%w", err)
	}

	// 先查询链接信息用于结果卡片展示
	link, err := c.store.GetLink(ctx, linkID)
	if err != nil {
		return "", "", fmt.Errorf("链接未找到:%w", err)
	}

	var responseText string
	switch actionVal.Action {
	case "link_approve":
		status := model.LinkStatusApproved
		if err := c.store.UpdateLink(ctx, &model.UpdateLink{Id: linkID, Status: &status}); err != nil {
			return "", "", fmt.Errorf("审核失败:%w", err)
		}
		responseText = "已通过审核"
	case "link_reject":
		if err := c.store.DeleteLink(ctx, linkID); err != nil {
			return "", "", fmt.Errorf("删除失败:%w", err)
		}
		responseText = "已拒绝并删除"
	default:
		return "", "", nil
	}

	msg := LinkMessage{
		Name:   link.Name,
		URL:    link.Url,
		Desc:   link.Desc,
		Result: responseText,
	}
	return c.BuildResultCard(msg), responseText, nil
}

// linkCardTemplate 友情链接审核卡片模板，包含"通过"和"拒绝"两个回调按钮
const linkCardTemplate = `{
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
                "content": {{.Name | json}},
                "text_align": "left",
                "text_size": "normal",
                "margin": "0px 0px 0px 0px",
                "icon": {
                    "tag": "standard_icon",
                    "token": "pen_outlined",
                    "color": "blue"
                }
            },
            {
                "tag": "markdown",
                "content": {{.URL | json}},
                "text_align": "left",
                "text_size": "notation",
                "margin": "0px 0px 0px 0px",
                "icon": {
                    "tag": "standard_icon",
                    "token": "sharelink_outlined",
                    "color": "blue"
                }
            },
            {
                "tag": "markdown",
                "content": {{.Desc | json}},
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
                                    "content": "通过"
                                },
                                "type": "primary_filled",
                                "width": "fill",
                                "size": "medium",
                                "behaviors": [
                                    {
                                        "type": "callback",
                                        "value": {
                                            "action": "link_approve",
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
                                    "content": "拒绝"
                                },
                                "type": "default",
                                "width": "fill",
                                "size": "medium",
                                "behaviors": [
                                    {
                                        "type": "callback",
                                        "value": {
                                            "action": "link_reject",
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
            "content": "友情链接审核"
        },
        "subtitle": {
            "tag": "plain_text",
            "content": ""
        },
        "template": "wathet",
        "padding": "12px 12px 12px 12px"
    }
}`

// linkResultCardTemplate 友情链接审核结果卡片模板，展示审核结果文本
const linkResultCardTemplate = `{
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
                "content": {{.Name | json}},
                "text_align": "left",
                "text_size": "normal",
                "margin": "0px 0px 0px 0px",
                "icon": {
                    "tag": "standard_icon",
                    "token": "pen_outlined",
                    "color": "blue"
                }
            },
            {
                "tag": "markdown",
                "content": {{.URL | json}},
                "text_align": "left",
                "text_size": "notation",
                "margin": "0px 0px 0px 0px",
                "icon": {
                    "tag": "standard_icon",
                    "token": "sharelink_outlined",
                    "color": "blue"
                }
            },
            {
                "tag": "markdown",
                "content": {{.Desc | json}},
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
            "content": "友情链接审核"
        },
        "subtitle": {
            "tag": "plain_text",
            "content": ""
        },
        "template": "wathet",
        "padding": "12px 12px 12px 12px"
    }
}`
