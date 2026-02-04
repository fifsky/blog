// Package aiutil provides utility functions for AI model configuration
package aiutil

import (
	"strings"

	"github.com/openai/openai-go/v3"
)

// ConfigureModelParams configures model-specific parameters for ChatCompletionNewParams.
// For doubao models, it disables thinking mode.
func ConfigureModelParams(req *openai.ChatCompletionNewParams, model string) {
	if strings.HasPrefix(model, "doubao") || strings.HasPrefix(model, "kimi") {
		req.SetExtraFields(map[string]any{
			"thinking": map[string]any{
				"type": "disabled",
			},
		})
	}
}
