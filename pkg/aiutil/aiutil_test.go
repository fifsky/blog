package aiutil

import (
	"encoding/json"
	"testing"

	"github.com/openai/openai-go/v3"
)

func TestConfigureModelParamsEnablesDeepSeekThinking(t *testing.T) {
	req := openai.ChatCompletionNewParams{}

	ConfigureModelParams(&req, "deepseek-v4-pro")

	if req.ReasoningEffort != "high" {
		t.Fatalf("reasoning effort = %q, want high", req.ReasoningEffort)
	}
	raw, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	var got struct {
		Thinking struct {
			Type string `json:"type"`
		} `json:"thinking"`
	}
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal request: %v", err)
	}
	if got.Thinking.Type != "enabled" {
		t.Fatalf("thinking.type = %q, want enabled; raw=%s", got.Thinking.Type, raw)
	}
}
