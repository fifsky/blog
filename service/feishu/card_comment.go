package feishu

import "text/template"

// CommentMessage 评论通知卡片消息
type CommentMessage struct {
	Name      string // 评论者昵称
	Content   string // 评论内容
	PostTitle string // 文章标题
	PostURL   string // 文章链接
	Time      string // 评论时间
}

// CommentCard 评论通知卡片处理器，仅通知无回调操作
type CommentCard struct {
	tplBuilder
}

// NewCommentCard 创建评论通知卡片处理器
func NewCommentCard() *CommentCard {
	return &CommentCard{
		tplBuilder: tplBuilder{
			cardTpl:   template.Must(template.New("comment").Funcs(tplFuncs).Parse(commentCardTemplate)),
			resultTpl: nil, // 无结果卡片
		},
	}
}

// BuildCard 构建评论通知卡片
func (c *CommentCard) BuildCard(msg CommentMessage) string { return c.execCard(msg) }

// commentCardTemplate 评论通知卡片模板，纯展示无按钮
const commentCardTemplate = `{
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
                "content": "**收到新评论**",
                "text_align": "left",
                "text_size": "heading",
                "margin": "0px 0px 0px 0px"
            },
            {
                "tag": "markdown",
                "content": {{.PostTitle | json}},
                "text_align": "left",
                "text_size": "normal",
                "margin": "0px 0px 0px 0px",
                "icon": {
                    "tag": "standard_icon",
                    "token": "doc_colorful",
                    "color": "blue"
                }
            },
            {
                "tag": "markdown",
                "content": {{.Name | json}},
                "text_align": "left",
                "text_size": "normal",
                "margin": "0px 0px 0px 0px",
                "icon": {
                    "tag": "standard_icon",
                    "token": "member_outlined",
                    "color": "blue"
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
            "content": "博客评论通知"
        },
        "subtitle": {
            "tag": "plain_text",
            "content": ""
        },
        "template": "turquoise",
        "padding": "12px 12px 12px 12px"
    }
}`
