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
	Name         string          `json:"name"`         // Unique name (mcpKey:originalName)
	OriginalName string          `json:"originalName"` // Original tool name from MCP
	MCPKey       string          `json:"mcpKey"`       // MCP client key
	Description  string          `json:"description"`
	InputSchema  json.RawMessage `json:"inputSchema"`
}

// Client wraps the MCP SDK client for MCP functionality
type Client struct {
	name   string
	url    string
	token  string
	client *mcp.Client

	mu      sync.Mutex
	session *mcp.ClientSession
	tools   []Tool
}

// NewClient creates a new MCP client
func NewClient(name, url, token string) *Client {
	return &Client{
		name:  name,
		url:   url,
		token: token,
		client: mcp.NewClient(&mcp.Implementation{
			Name:    "blog-ai",
			Version: "v1.0.0",
		}, nil),
	}
}

// Name returns the client name
func (c *Client) Name() string {
	return c.name
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
		return fmt.Errorf("failed to connect to MCP server %s: %w", c.name, err)
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
		return nil, fmt.Errorf("failed to list tools from %s: %w", c.name, err)
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

// Manager manages multiple MCP clients and provides aggregated tool access
type Manager struct {
	clients         map[string]*Client
	toolToMCP       map[string]string // maps tool name to MCP client key
	mcpDisplayNames map[string]string // maps MCP client key to display name
	mu              sync.RWMutex
}

// NewManager creates a new MCP manager with multiple clients
func NewManager() *Manager {
	return &Manager{
		clients:         make(map[string]*Client),
		toolToMCP:       make(map[string]string),
		mcpDisplayNames: make(map[string]string),
	}
}

// AddClient adds an MCP client to the manager
func (m *Manager) AddClient(key, displayName, url, token string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[key] = NewClient(key, url, token)
	m.mcpDisplayNames[key] = displayName
}

// ListAllTools returns all tools from all MCP clients
func (m *Manager) ListAllTools(ctx context.Context) ([]Tool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var allTools []Tool
	m.toolToMCP = make(map[string]string) // reset mapping

	for mcpKey, client := range m.clients {
		tools, err := client.ListTools(ctx)
		if err != nil {
			// Log error but continue with other clients
			continue
		}
		for _, tool := range tools {
			// Create unique tool name: mcpKey:originalName
			uniqueName := mcpKey + ":" + tool.Name
			m.toolToMCP[uniqueName] = mcpKey
			allTools = append(allTools, Tool{
				Name:         uniqueName,
				OriginalName: tool.Name,
				MCPKey:       mcpKey,
				Description:  tool.Description,
				InputSchema:  tool.InputSchema,
			})
		}
	}

	return allTools, nil
}

// CallTool calls a tool, routing to the correct MCP client
// toolName should be in format "mcpKey:originalToolName"
func (m *Manager) CallTool(ctx context.Context, toolName string, arguments map[string]any) (string, error) {
	m.mu.RLock()
	mcpKey, ok := m.toolToMCP[toolName]
	client := m.clients[mcpKey]
	m.mu.RUnlock()

	if !ok || client == nil {
		return "", fmt.Errorf("tool %s not found in any MCP client", toolName)
	}

	// Extract original tool name (after the colon)
	originalName := toolName
	if idx := len(mcpKey) + 1; idx < len(toolName) {
		originalName = toolName[idx:]
	}

	return client.CallTool(ctx, originalName, arguments)
}

// GetMCPDisplayName returns the display name for a tool's MCP client
func (m *Manager) GetMCPDisplayName(toolName string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	mcpKey := m.toolToMCP[toolName]
	if displayName, ok := m.mcpDisplayNames[mcpKey]; ok {
		return displayName
	}
	return mcpKey
}

// Close closes all MCP client connections
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, client := range m.clients {
		_ = client.Close()
	}
	return nil
}

// HasClients returns true if there are any MCP clients configured
func (m *Manager) HasClients() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.clients) > 0
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
