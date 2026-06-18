# weixin-clawbot

用于 OpenClaw 微信接入的 Go 包，提供微信扫码登录、消息收发、长轮询监听以及媒体上传下载能力。

## 项目简介

**weixin-clawbot** 是一个 Go 语言库，专为通过 [iLink](https://ilinkai.weixin.qq.com) Bot 平台接入微信（Weixin/WeChat）而设计，是 OpenClaw 微信插件的底层客户端。

它解决的核心问题是：如何用 Go 代码以机器人身份登录微信、接收消息并自动回复。整个流程包括：

1. **扫码登录**：返回二维码内容给页面或其他 UI 展示，用户用手机扫码授权后，库会持久化登录凭证，后续无需重复登录。
2. **消息监听**：通过长轮询（Long Polling）实时接收来自微信好友或群组的消息。
3. **消息发送**：向指定用户或群组发送文本、图片、视频、文件等多种类型的消息。
4. **媒体处理**：支持将媒体文件上传至 CDN（含 AES-ECB 加密），以及下载入站媒体文件到本地。

如果你正在构建一个基于微信的聊天机器人、消息自动化系统或客服机器人，这个库提供了所需的底层能力。

接口协议详见 [API 文档](doc/protocol.md)

## 功能

- 微信扫码登录，并可将账号信息保存到本地
- 支持长轮询监听与 `get_updates_buf` 持久化
- 提供文本、图片、视频、文件发送辅助函数
- 提供带 AES-ECB 处理的 CDN 上传下载能力
- 提供入站媒体落盘与消息转换工具

## 快速开始

### 1. 登录

```go
package main

import (
	"context"
	"log"

	"app/pkg/clawbot"
)

func main() {
	client := clawbot.NewClient(clawbot.Options{})

	session, err := client.StartLogin(context.Background(), "")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("show this QR content in your page: %s", session.QRContent)

	account, err := client.WaitLogin(context.Background(), session, clawbot.WaitOptions{SaveDir: ".weixin-accounts"})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account=%s token=%s", account.AccountID, account.BotToken)
}
```

### 2. 发送文本消息

```go
ctx := context.Background()

sender := clawbot.NewSender(clawbot.SenderOptions{
	BaseURL:      "https://ilinkai.weixin.qq.com",
	Token:        "YOUR_BOT_TOKEN",
	Timeout:      15 * time.Second,
})
conversation := sender.Conversation(clawbot.Target{
	ToUserID:     "user@im.wechat",
	ContextToken: "YOUR_CONTEXT_TOKEN",
})

clientID, err := conversation.SendText(ctx, "hello from bot")
if err != nil {
	log.Fatal(err)
}

log.Println("sent:", clientID)
```

### 3. 监听消息

```go
api := clawbot.NewAPIClient(clawbot.APIOptions{
	BaseURL: "https://ilinkai.weixin.qq.com",
	Token:   "YOUR_BOT_TOKEN",
})

err := clawbot.Listen(context.Background(), clawbot.ListenOptions{
	API:         api,
	AccountID:   "bot@im.bot",
	SyncBufPath: clawbot.SyncBufFilePath(clawbot.ResolveStateDir(), "bot@im.bot"),
	OnMessages: func(ctx context.Context, messages []clawbot.WeixinMessage) error {
		for _, msg := range messages {
			log.Printf("from=%s body=%s", msg.FromUserID, clawbot.BodyFromItemList(msg.ItemList))
		}
		return nil
	},
})
if err != nil {
	log.Fatal(err)
}
```

## 主要类型

- `Client`：扫码登录流程
- `APIClient`：ilink bot API 封装
- `Sender`：可复用的消息发送器
- `Conversation`：绑定单个 `ToUserID` + `ContextToken` 的会话发送器
- `Target`：发送目标
- `ListenOptions`：长轮询监听配置
- `UploadedFileInfo`：CDN 上传结果

## 说明

- 账号文件名会做 base64url 编码，避免路径字符不安全。
- `Target.ContextToken` 是发送消息时的必填字段。
