package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Tool represents a tool available from the MCP server
type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

// Client wraps the MCP SDK client for web search functionality
type Client struct {
	url    string
	token  string
	client *mcp.Client

	mu      sync.Mutex
	session *mcp.ClientSession
	tools   []Tool
}

// NewClient creates a new MCP client for web search
func NewClient(url, token string) *Client {
	return &Client{
		url:   url,
		token: token,
		client: mcp.NewClient(&mcp.Implementation{
			Name:    "blog-ai",
			Version: "v1.0.0",
		}, nil),
	}
}

// Connect establishes a connection to the MCP server
func (c *Client) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Close existing session if any
	if c.session != nil {
		_ = c.session.Close()
		c.session = nil
	}

	// Create HTTP client with auth header
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &authTransport{
			token: c.token,
			rt:    http.DefaultTransport,
		},
	}

	transport := &mcp.StreamableClientTransport{
		Endpoint:   c.url,
		HTTPClient: httpClient,
	}

	session, err := c.client.Connect(ctx, transport, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to MCP server: %w", err)
	}

	c.session = session
	return nil
}

// ensureSession ensures there is an active session
func (c *Client) ensureSession(ctx context.Context) (*mcp.ClientSession, error) {
	c.mu.Lock()
	session := c.session
	c.mu.Unlock()

	if session == nil {
		if err := c.Connect(ctx); err != nil {
			return nil, err
		}
		c.mu.Lock()
		session = c.session
		c.mu.Unlock()
	}

	return session, nil
}

// ListTools returns the list of available tools from the MCP server
func (c *Client) ListTools(ctx context.Context) ([]Tool, error) {
	session, err := c.ensureSession(ctx)
	if err != nil {
		return nil, err
	}

	result, err := session.ListTools(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}

	tools := make([]Tool, 0, len(result.Tools))
	for _, t := range result.Tools {
		schemaBytes, _ := json.Marshal(t.InputSchema)
		tools = append(tools, Tool{
			Name:        t.Name,
			Description: t.Description,
			InputSchema: schemaBytes,
		})
	}

	c.mu.Lock()
	c.tools = tools
	c.mu.Unlock()

	return tools, nil
}

// CallTool calls a tool on the MCP server
func (c *Client) CallTool(ctx context.Context, name string, arguments map[string]any) (string, error) {
	session, err := c.ensureSession(ctx)
	if err != nil {
		return "", err
	}

	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      name,
		Arguments: arguments,
	})
	if err != nil {
		return "", fmt.Errorf("tool call failed: %w", err)
	}

	if result.IsError {
		return "", fmt.Errorf("tool returned error")
	}

	// Extract text content from result
	var content string
	for _, c := range result.Content {
		if textContent, ok := c.(*mcp.TextContent); ok {
			content += textContent.Text
		}
	}

	return content, nil
}

// Close closes the MCP client connection
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.session != nil {
		err := c.session.Close()
		c.session = nil
		return err
	}
	return nil
}

// authTransport adds authorization header to all requests
type authTransport struct {
	token string
	rt    http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.token != "" {
		req.Header.Set("Authorization", "Bearer "+t.token)
	}
	return t.rt.RoundTrip(req)
}
