package feishu

import (
	"bytes"
	"context"
	"encoding/json"
	"text/template"

	"app/config"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// FeishuSender 通过飞书机器人发送卡片消息，只负责发送 JSON 字符串
type FeishuSender struct {
	client *lark.Client
	userID string
}

// NewFeishuSender 创建飞书发送器，飞书配置缺失时返回 nil
func NewFeishuSender(conf config.FeishuConf) *FeishuSender {
	if conf.Appid == "" || conf.AppSecret == "" || conf.UserID == "" {
		return nil
	}
	return &FeishuSender{
		client: lark.NewClient(conf.Appid, conf.AppSecret),
		userID: conf.UserID,
	}
}

// Send 发送卡片 JSON 到飞书
func (f *FeishuSender) Send(ctx context.Context, cardJSON string) error {
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType("open_id").
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(f.userID).
			MsgType("interactive").
			Content(cardJSON).
			Build()).
		Build()

	resp, err := f.client.Im.V1.Message.Create(ctx, req)
	if err != nil {
		return err
	}
	if !resp.Success() {
		return &FeishuError{Code: resp.Code, Msg: resp.Msg}
	}
	return nil
}

// FeishuError 飞书 API 错误
type FeishuError struct {
	Code int
	Msg  string
}

func (e *FeishuError) Error() string { return e.Msg }

// CardHandler 卡片回调处理器接口，每种业务卡片实现并注册到 CardRegistry。
// 卡片构建（BuildCard/BuildResultCard）是具体方法，不在接口上，
// 各卡片可定义自己的消息结构体，调用方直接调用具体卡片的构建方法。
type CardHandler interface {
	// Actions 返回该处理器负责的回调 action 列表
	Actions() []string
	// Handle 处理回调动作，内部构建结果卡片，返回卡片 JSON 和结果文本
	Handle(ctx context.Context, action, token string) (cardJSON string, resultText string, err error)
}

// CardRegistry 卡片注册表，按回调 action 分发到对应处理器
type CardRegistry struct {
	handlers map[string]CardHandler
}

// NewCardRegistry 创建卡片注册表
func NewCardRegistry() *CardRegistry {
	return &CardRegistry{handlers: make(map[string]CardHandler)}
}

// Register 注册卡片处理器，自动绑定回调 action
func (r *CardRegistry) Register(card CardHandler) {
	for _, a := range card.Actions() {
		r.handlers[a] = card
	}
}

// Handle 处理回调动作，返回结果卡片 JSON 和结果文本。
// 未注册的 action 返回空字符串，调用方可据此跳过
func (r *CardRegistry) Handle(ctx context.Context, action, token string) (cardJSON string, resultText string, err error) {
	handler, ok := r.handlers[action]
	if !ok {
		return "", "", nil
	}
	return handler.Handle(ctx, action, token)
}

// tplFuncs 模板公共函数，用于 JSON 序列化并自动转义特殊字符
var tplFuncs = template.FuncMap{
	"json": func(v any) (string, error) {
		b, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		return string(b), nil
	},
}

// tplBuilder 封装模板执行的通用逻辑，各卡片处理器嵌入复用。
// 接受 any 类型，各卡片传入自己的消息结构体
type tplBuilder struct {
	cardTpl   *template.Template
	resultTpl *template.Template
}

func (t *tplBuilder) execCard(msg any) string {
	var buf bytes.Buffer
	_ = t.cardTpl.Execute(&buf, msg)
	return buf.String()
}

func (t *tplBuilder) execResult(msg any) string {
	var buf bytes.Buffer
	_ = t.resultTpl.Execute(&buf, msg)
	return buf.String()
}
