package doubao

// ChatRequest 代表聊天完成接口的请求体
type ChatRequest struct {
	// Model 指定使用的模型 ID
	Model string `json:"model"`

	// Tools 定义模型可以调用的工具列表
	Tools []Tool `json:"tools,omitempty"`

	// MaxToolCalls 限制模型单次交互中调用工具的最大次数
	MaxToolCalls int `json:"max_tool_calls,omitempty"`

	// Thinking 配置模型的思考过程（如果支持）
	Thinking *Thinking `json:"thinking,omitempty"`

	// Input 包含对话的历史消息列表
	Input []Message `json:"input"`
}

// Tool 定义一个工具
type Tool struct {
	// Type 工具类型，例如 "web_search"
	Type string `json:"type"`

	// MaxKeyword 仅用于 web_search，限制搜索关键词数量
	MaxKeyword int `json:"max_keyword,omitempty"`

	// Limit 仅用于 web_search，限制搜索结果数量
	Limit int `json:"limit,omitempty"`
}

// Thinking 定义思考配置
type Thinking struct {
	// Type 思考类型，例如 "disabled"
	Type string `json:"type"`
}

// Message 代表一条对话消息
type Message struct {
	// Role 消息发送者的角色，例如 "system", "user", "assistant"
	Role string `json:"role"`

	// Content 消息的内容列表
	Content []MessageContent `json:"content"`
}

// MessageContent 代表消息的具体内容
type MessageContent struct {
	// Type 内容类型，例如 "input_text"
	Type string `json:"type"`

	// Text 具体的文本内容
	Text string `json:"text"`
}

// ChatResponse 代表聊天完成接口的响应体
type ChatResponse struct {
	Output []Output `json:"output"`
}

type Output struct {
	Type    string    `json:"type"`
	Content []Content `json:"content"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
