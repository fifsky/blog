package promptutil

import "os"

// ParsePrompt 解析prompt模板，将args中的值替换到模板中
// prompt 模板字符串，例如："你好，${name}！"
// args 模板参数，例如：map[string]string{"name": "张三"}
// 返回解析后的字符串，例如："你好，张三！"
func ParsePrompt(prompt string, args map[string]string) string {
	return os.Expand(prompt, func(key string) string {
		return args[key]
	})
}
