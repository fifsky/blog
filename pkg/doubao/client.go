package doubao

import (
	"net/http"
	"time"
)

const (
	DefaultBaseURL = "https://ark.cn-beijing.volces.com/api/v3"
)

// Client is the client for the Doubao API
type Client struct {
	// Base URL of the API
	BaseURL string

	// Bearer Token used for authentication
	APIKey string

	// HTTP client used to send requests
	HTTPClient *http.Client
}

// NewClient creates a new Client instance
func NewClient(apiKey string) *Client {
	return &Client{
		BaseURL: DefaultBaseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}
