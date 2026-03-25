package tool

import (
	"context"
	"encoding/json"
)

// Handler consumes tool arguments returned by the LLM (serialized as JSON string).
type Handler interface {
	Handle(ctx context.Context, arguments string) (string, error)
}

// HandleFunc adapts a plain function to a Handler.
type HandleFunc func(ctx context.Context, arguments string) (string, error)

func (f HandleFunc) Handle(ctx context.Context, arguments string) (string, error) {
	return f(ctx, arguments)
}

// Tool defines the interface for a tool that can be used in a system.
type Tool interface {
	Name() string
	Description() string
	InputSchema() json.RawMessage
	Handler
}

// baseTool is a simple implementation of Tool.
type baseTool struct {
	name        string
	description string
	inputSchema json.RawMessage
	handler     Handler
}

func (b *baseTool) Name() string {
	return b.name
}

func (b *baseTool) Description() string {
	return b.description
}

func (b *baseTool) InputSchema() json.RawMessage {
	return b.inputSchema
}

func (b *baseTool) Handle(ctx context.Context, arguments string) (string, error) {
	return b.handler.Handle(ctx, arguments)
}

// NewTool creates a new Tool with the given name, description, schema, and handler.
func NewTool(name, description string, inputSchema json.RawMessage, handler Handler) Tool {
	return &baseTool{
		name:        name,
		description: description,
		inputSchema: inputSchema,
		handler:     handler,
	}
}

// JSONAdapter adapts a typed function to a HandleFunc.
func JSONAdapter[I, O any](handle func(context.Context, I) (O, error)) HandleFunc {
	return func(ctx context.Context, input string) (string, error) {
		var req I
		if err := json.Unmarshal([]byte(input), &req); err != nil {
			return "", err
		}
		res, err := handle(ctx, req)
		if err != nil {
			return "", err
		}
		b, err := json.Marshal(res)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
}
