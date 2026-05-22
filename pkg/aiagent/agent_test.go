package aiagent

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	mcpclient "app/pkg/mcp"

	"github.com/openai/openai-go/v3"
)

type fakeToolProvider struct {
	tools       []mcpclient.Tool
	result      string
	displayName string
	calls       []string
	args        []map[string]any
}

func (f *fakeToolProvider) HasClients() bool {
	return true
}

func (f *fakeToolProvider) ListAllTools(context.Context) ([]mcpclient.Tool, error) {
	return f.tools, nil
}

func (f *fakeToolProvider) CallTool(_ context.Context, name string, arguments map[string]any) (string, error) {
	f.calls = append(f.calls, name)
	f.args = append(f.args, arguments)
	return f.result, nil
}

func (f *fakeToolProvider) GetMCPDisplayName(string) string {
	return f.displayName
}

type fakeStreamFactory struct {
	requests []openai.ChatCompletionNewParams
	streams  []*fakeStream
}

func (f *fakeStreamFactory) NewStreaming(_ context.Context, req openai.ChatCompletionNewParams) (chatStream, error) {
	f.requests = append(f.requests, req)
	if len(f.streams) == 0 {
		return nil, errors.New("no stream configured")
	}
	stream := f.streams[0]
	f.streams = f.streams[1:]
	return stream, nil
}

type fakeStream struct {
	chunks []openai.ChatCompletionChunk
	index  int
	err    error
}

func (f *fakeStream) Next() bool {
	if f.index >= len(f.chunks) {
		return false
	}
	f.index++
	return true
}

func (f *fakeStream) Current() openai.ChatCompletionChunk {
	return f.chunks[f.index-1]
}

func (f *fakeStream) Err() error {
	return f.err
}

func TestBuildTools(t *testing.T) {
	tests := []struct {
		name       string
		schema     json.RawMessage
		wantSchema map[string]any
	}{
		{
			name:       "valid schema",
			schema:     json.RawMessage(`{"type":"object","properties":{"title":{"type":"string"}}}`),
			wantSchema: map[string]any{"type": "object", "properties": map[string]any{"title": map[string]any{"type": "string"}}},
		},
		{
			name:       "empty schema falls back to object",
			schema:     nil,
			wantSchema: map[string]any{"type": "object"},
		},
		{
			name:       "invalid schema falls back to object",
			schema:     json.RawMessage(`{`),
			wantSchema: map[string]any{"type": "object"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := &Agent{
				tools: &fakeToolProvider{
					tools: []mcpclient.Tool{{
						Name:        "blog:remind",
						Description: "create remind",
						InputSchema: tt.schema,
					}},
				},
			}

			tools := agent.buildTools(context.Background())
			if len(tools) != 1 {
				t.Fatalf("len(tools) = %d, want 1", len(tools))
			}

			got := tools[0].OfFunction.Function.Parameters
			if !reflect.DeepEqual(map[string]any(got), tt.wantSchema) {
				t.Fatalf("schema = %#v, want %#v", got, tt.wantSchema)
			}
		})
	}
}

func TestMemory(t *testing.T) {
	now := time.Date(2026, 5, 22, 9, 0, 0, 0, time.Local)
	memory := NewMemory(time.Hour, 4)
	memory.now = func() time.Time { return now }

	memory.Save("sender", []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage("one"),
		openai.AssistantMessage("two"),
		openai.UserMessage("three"),
		openai.AssistantMessage("four"),
		openai.UserMessage("five"),
		openai.AssistantMessage("six"),
	})

	got := memory.Get("sender")
	if len(got) != 4 {
		t.Fatalf("len(got) = %d, want 4", len(got))
	}
	if got[0].OfUser == nil || got[0].OfUser.Content.OfString.Value != "three" {
		t.Fatalf("first retained message = %#v, want user three", got[0])
	}

	got[0] = openai.UserMessage("mutated")
	gotAgain := memory.Get("sender")
	if gotAgain[0].OfUser.Content.OfString.Value != "three" {
		t.Fatalf("memory returned mutable backing slice")
	}

	memory.now = func() time.Time { return now.Add(time.Hour + time.Second) }
	if expired := memory.Get("sender"); expired != nil {
		t.Fatalf("expired messages = %#v, want nil", expired)
	}
}

func TestAgentRunStreamsContentAndExecutesTools(t *testing.T) {
	toolProvider := &fakeToolProvider{
		tools: []mcpclient.Tool{{
			Name:        "blog:weather",
			Description: "query weather",
			InputSchema: json.RawMessage(`{"type":"object"}`),
		}},
		result:      "sunny",
		displayName: "天气",
	}
	streamFactory := &fakeStreamFactory{
		streams: []*fakeStream{
			{
				chunks: []openai.ChatCompletionChunk{
					chunk(`{"choices":[{"index":0,"delta":{"role":"assistant","tool_calls":[{"index":0,"id":"call_1","type":"function","function":{"name":"blog:weather","arguments":"{\"city\":\"上海\"}"}}]},"finish_reason":"tool_calls"}]}`),
				},
			},
			{
				chunks: []openai.ChatCompletionChunk{
					chunk(`{"choices":[{"index":0,"delta":{"role":"assistant","content":"今日晴朗"},"finish_reason":null}]}`),
					chunk(`{"choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}`),
				},
			},
		},
	}
	agent := &Agent{
		clientProvider: func(context.Context) (openai.Client, string, error) {
			return openai.Client{}, "test-model", nil
		},
		streamFactory: streamFactory,
		tools:         toolProvider,
	}

	var events []string
	var content string
	result, err := agent.Run(context.Background(), Request{
		SystemPrompt: "system",
		Messages:     []openai.ChatCompletionMessageParamUnion{openai.UserMessage("今天天气？")},
		UseTools:     true,
	}, EventHandler{
		OnContent: func(_ context.Context, delta string) error {
			events = append(events, "content:"+delta)
			content += delta
			return nil
		},
		OnToolStart: func(_ context.Context, event ToolEvent) error {
			events = append(events, "start:"+event.ID+":"+event.Name+":"+event.MCPName+":"+event.Arguments)
			return nil
		},
		OnToolEnd: func(_ context.Context, event ToolEvent) error {
			events = append(events, "end:"+event.ID+":"+event.Result)
			return nil
		},
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	wantEvents := []string{
		`start:call_1:blog:weather:天气:{"city":"上海"}`,
		"end:call_1:sunny",
		"content:今日晴朗",
	}
	if !reflect.DeepEqual(events, wantEvents) {
		t.Fatalf("events = %#v, want %#v", events, wantEvents)
	}
	if content != "今日晴朗" || result.Content != content {
		t.Fatalf("content = %q, result = %#v", content, result)
	}
	if !reflect.DeepEqual(toolProvider.calls, []string{"blog:weather"}) {
		t.Fatalf("tool calls = %#v, want blog:weather", toolProvider.calls)
	}
	if got := toolProvider.args[0]["city"]; got != "上海" {
		t.Fatalf("tool arg city = %#v, want 上海", got)
	}
	if len(streamFactory.requests) != 2 {
		t.Fatalf("stream requests = %d, want 2", len(streamFactory.requests))
	}
	if len(streamFactory.requests[0].Tools) != 1 {
		t.Fatalf("first request tools = %d, want 1", len(streamFactory.requests[0].Tools))
	}
	if len(streamFactory.requests[1].Messages) != 4 {
		t.Fatalf("second request messages = %d, want system/user/assistant/tool", len(streamFactory.requests[1].Messages))
	}
}

func TestAgentRunReturnsReasoningContentAfterToolCalls(t *testing.T) {
	toolProvider := &fakeToolProvider{
		tools: []mcpclient.Tool{{
			Name:        "blog:weather",
			Description: "query weather",
			InputSchema: json.RawMessage(`{"type":"object"}`),
		}},
		result: "sunny",
	}
	streamFactory := &fakeStreamFactory{
		streams: []*fakeStream{
			{
				chunks: []openai.ChatCompletionChunk{
					chunk(`{"choices":[{"index":0,"delta":{"role":"assistant","reasoning_content":"需要先查天气。"},"finish_reason":null}]}`),
					chunk(`{"choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"call_1","type":"function","function":{"name":"blog:weather","arguments":"{\"city\":\"上海\"}"}}]},"finish_reason":"tool_calls"}]}`),
				},
			},
			{
				chunks: []openai.ChatCompletionChunk{
					chunk(`{"choices":[{"index":0,"delta":{"role":"assistant","reasoning_content":"工具结果可用。","content":"今日晴朗"},"finish_reason":null}]}`),
					chunk(`{"choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}`),
				},
			},
		},
	}
	agent := &Agent{
		clientProvider: func(context.Context) (openai.Client, string, error) {
			return openai.Client{}, "deepseek-v4-pro", nil
		},
		streamFactory: streamFactory,
		tools:         toolProvider,
	}

	result, err := agent.Run(context.Background(), Request{
		Messages: []openai.ChatCompletionMessageParamUnion{openai.UserMessage("今天天气？")},
		UseTools: true,
	}, EventHandler{})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if len(streamFactory.requests[1].Messages) != 3 {
		t.Fatalf("second request messages = %d, want user/assistant/tool", len(streamFactory.requests[1].Messages))
	}
	assertMessageReasoningContent(t, streamFactory.requests[1].Messages[1], "需要先查天气。")
	if len(result.Messages) != 3 {
		t.Fatalf("result messages = %d, want assistant/tool/assistant", len(result.Messages))
	}
	assertMessageReasoningContent(t, result.Messages[0], "需要先查天气。")
	assertMessageReasoningContent(t, result.Messages[2], "工具结果可用。")
}

func TestAssistantMessageReasoningContentSurvivesJSONRoundTrip(t *testing.T) {
	message := openai.ChatCompletionMessage{
		Content: "查好了",
	}
	msg := assistantMessageWithReasoning(message, "需要保留。")

	raw, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal message: %v", err)
	}
	decoded, err := DecodeMessageParam(raw)
	if err != nil {
		t.Fatalf("unmarshal message: %v", err)
	}

	assertMessageReasoningContent(t, decoded, "需要保留。")
}

func assertMessageReasoningContent(t *testing.T, msg openai.ChatCompletionMessageParamUnion, want string) {
	t.Helper()

	raw, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal message: %v", err)
	}
	var got struct {
		ReasoningContent string `json:"reasoning_content"`
	}
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal message: %v", err)
	}
	if got.ReasoningContent != want {
		t.Fatalf("reasoning_content = %q, want %q; raw=%s", got.ReasoningContent, want, raw)
	}
}

func chunk(raw string) openai.ChatCompletionChunk {
	var c openai.ChatCompletionChunk
	if err := json.Unmarshal([]byte(raw), &c); err != nil {
		panic(err)
	}
	return c
}
