package feishu

import (
	"context"
	"fmt"
	"strconv"
	"text/template"
	"time"

	"app/config"
	"app/pkg/aesutil"
	"app/pkg/remindutil"
	"app/store"
	"app/store/model"

	"github.com/samber/lo"
)

// RemindActionValue 提醒卡片回调的 actionValue 结构体
type RemindActionValue struct {
	Action string `json:"action"` // 回调动作：remind_completed / remind_later
	Token  string `json:"token"`  // 加密的业务 ID
}

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
	store *store.Store
	conf  *config.Config
}

// NewRemindCard 创建提醒卡片处理器
func NewRemindCard(store *store.Store, conf *config.Config) *RemindCard {
	return &RemindCard{
		tplBuilder: tplBuilder{
			cardTpl:   template.Must(template.New("remind").Funcs(tplFuncs).Parse(remindCardTemplate)),
			resultTpl: template.Must(template.New("remindResult").Funcs(tplFuncs).Parse(remindResultCardTemplate)),
		},
		store: store,
		conf:  conf,
	}
}

func (c *RemindCard) Actions() []string { return []string{"remind_completed", "remind_later"} }

// BuildCard 构建提醒卡片
func (c *RemindCard) BuildCard(msg RemindMessage) string { return c.execCard(msg) }

// BuildResultCard 构建提醒结果卡片
func (c *RemindCard) BuildResultCard(msg RemindMessage) string { return c.execResult(msg) }

// Handle 处理提醒卡片回调，返回结果卡片 JSON 和结果文本
func (c *RemindCard) Handle(ctx context.Context, actionValue map[string]any) (string, string, error) {
	actionVal, err := parseActionValue[RemindActionValue](actionValue)
	if err != nil {
		return "", "", fmt.Errorf("解析actionValue失败: %w", err)
	}

	// 解密 token 获取提醒 ID
	id, err := aesutil.AesDecode(c.conf.Common.TokenSecret, actionVal.Token)
	if err != nil {
		return "", "", fmt.Errorf("token错误: %w", err)
	}

	remind, err := c.store.GetRemind(ctx, lo.Must(strconv.Atoi(id)))
	if err != nil {
		return "", "", fmt.Errorf("记录未找到: %w", err)
	}

	var resultText string
	switch actionVal.Action {
	case "remind_completed":
		resultText, err = c.handleChange(ctx, remind)
	case "remind_later":
		resultText, err = c.handleDelay(ctx, remind)
	default:
		return "", "", nil
	}
	if err != nil {
		return "", "", fmt.Errorf("操作失败: %w", err)
	}

	// 重新查询获取更新后的 remind（nextTime 可能已变更）
	remind, err = c.store.GetRemind(ctx, lo.Must(strconv.Atoi(id)))
	if err != nil {
		return "", "", fmt.Errorf("记录未找到: %w", err)
	}

	msg := RemindMessage{
		Content: remind.Content,
		Time:    remind.NextTime.Format("2006-01-02 15:04"),
		Result:  resultText,
	}
	return c.BuildResultCard(msg), resultText, nil
}

// handleChange 标记完成：固定时间任务设为已完成，周期性任务恢复状态并计算下次时间
func (c *RemindCard) handleChange(ctx context.Context, remind *model.Remind) (string, error) {
	nextTime := remindutil.NextTimeFromRule(time.Now(), remind)

	// 判断是否是固定时间任务（cron 格式为 2006-01-02 15:04:xx）
	isFixedDate := len(remind.Cron) >= 10 && remind.Cron[4] == '-'

	if isFixedDate {
		if err := c.store.UpdateRemindStatus(ctx, remind.Id, 3); err != nil {
			return "", err
		}
		return "已确认完成", nil
	}

	// 周期性任务，恢复状态并更新下次时间
	if err := c.store.UpdateRemindStatus(ctx, remind.Id, 1); err != nil {
		return "", err
	}
	if err := c.store.UpdateRemindNextTime(ctx, remind.Id, nextTime); err != nil {
		return "", err
	}
	return "已确认收到提醒", nil
}

// handleDelay 延迟提醒：下次提醒时间推迟 10 分钟
func (c *RemindCard) handleDelay(ctx context.Context, remind *model.Remind) (string, error) {
	nextTime := time.Now().Add(10 * time.Minute)
	if err := c.store.UpdateRemindNextTime(ctx, remind.Id, nextTime); err != nil {
		return "", err
	}
	return "将在10分钟后再次提醒", nil
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
