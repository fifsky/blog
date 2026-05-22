// Package aiutil 提供 AI 模型请求参数配置工具。
package aiutil

import (
	"strings"

	"github.com/openai/openai-go/v3"
)

// ConfigureModelParams 按模型厂商配置 ChatCompletion 请求参数。
func ConfigureModelParams(req *openai.ChatCompletionNewParams, model string) {
	if strings.HasPrefix(model, "doubao") {
		req.ReasoningEffort = "minimal"
	}

	if strings.HasPrefix(model, "kimi") {
		req.SetExtraFields(map[string]any{
			"thinking": map[string]any{
				"type": "disabled",
			},
		})
	}

	if strings.HasPrefix(model, "deepseek") {
		req.ReasoningEffort = "high"
		req.SetExtraFields(map[string]any{
			"thinking": map[string]any{
				"type": "enabled",
			},
		})
	}
}
