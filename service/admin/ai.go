package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"app/config"
	"app/pkg/errors"
	"app/server/response"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type AI struct {
	conf   *config.Config
	client openai.Client
}

func NewAI(conf *config.Config) *AI {
	client := openai.NewClient(
		option.WithAPIKey(conf.Common.AIToken),
		option.WithBaseURL(conf.Common.AIEndpoint),
	)
	return &AI{
		conf:   conf,
		client: client,
	}
}

// ChatMessage represents a single message in the conversation
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents the incoming chat request with message history
type ChatRequest struct {
	Messages []ChatMessage `json:"messages"`
}

// Chat handles SSE streaming AI chat responses
func (a *AI) Chat(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, errors.BadRequest("INVALID_REQUEST", "Invalid request body"))
		return
	}

	if len(req.Messages) == 0 {
		response.Fail(w, errors.BadRequest("EMPTY_MESSAGES", "Messages cannot be empty"))
		return
	}

	// Validate that there's at least one user message with content
	hasContent := false
	for _, msg := range req.Messages {
		if strings.TrimSpace(msg.Content) != "" {
			hasContent = true
			break
		}
	}
	if !hasContent {
		response.Fail(w, errors.BadRequest("EMPTY_MESSAGE", "Message content cannot be empty"))
		return
	}

	// Set SSE headers early
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		response.Fail(w, errors.InternalServer("STREAMING_ERROR", "Streaming not supported"))
		return
	}

	prompt := fmt.Sprintf(`You are a helpful assistant. Respond in the same language as the user's message.
Current Time: %s
	`, time.Now().Format(time.DateTime))

	// Build OpenAI messages from request history
	openAIMessages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(prompt),
	}

	for _, msg := range req.Messages {
		if strings.TrimSpace(msg.Content) == "" {
			continue // skip empty messages
		}
		switch msg.Role {
		case "user":
			openAIMessages = append(openAIMessages, openai.UserMessage(msg.Content))
		case "assistant":
			openAIMessages = append(openAIMessages, openai.AssistantMessage(msg.Content))
		}
	}

	// Create streaming chat completion using OpenAI SDK v3
	aiReq := openai.ChatCompletionNewParams{
		Model:    a.conf.Common.AIModel,
		Messages: openAIMessages,
	}
	if strings.HasPrefix(a.conf.Common.AIModel, "doubao") {
		aiReq.SetExtraFields(map[string]any{
			"thinking": map[string]any{
				"type": "disabled",
			},
		})
	}

	stream := a.client.Chat.Completions.NewStreaming(r.Context(), aiReq)

	// Stream the response
	for stream.Next() {
		chunk := stream.Current()
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			content := chunk.Choices[0].Delta.Content
			// Escape newlines for SSE
			escapedContent := strings.ReplaceAll(content, "\n", "\\n")
			fmt.Fprintf(w, "data: %s\n\n", escapedContent)
			flusher.Flush()
		}
	}

	if err := stream.Err(); err != nil {
		// Check if context was cancelled
		if r.Context().Err() == context.Canceled {
			return
		}
		// Send error as SSE event
		fmt.Fprintf(w, "data: [ERROR] %s\n\n", err.Error())
		flusher.Flush()
		return
	}

	// Send done event
	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()
}
