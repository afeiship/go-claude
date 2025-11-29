package claude

import (
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

// Config defines the client configuration
type Config struct {
	APIKey    string        // Required
	BaseURL   string        // Optional, default "https://api.anthropic.com"
	Model     string        // Optional, default "claude-3-5-sonnet-20241022"
	MaxTokens int           // Optional, default 1024
	Timeout   time.Duration // Optional, default 60s
}

// Client is the Claude API client
type Client struct {
	config Config
	client *resty.Client
}

// Message represents a conversation message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type requestPayload struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []Message `json:"messages"`
}

type contentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type apiResponse struct {
	Content []contentBlock `json:"content"`
}

// NewClient creates a client from Config
func NewClient(cfg Config) (*Client, error) {
	// Use environment variables as fallback (optional, you can also disable this completely)
	if cfg.APIKey == "" {
		cfg.APIKey = os.Getenv("ANTHROPIC_AUTH_TOKEN")
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = os.Getenv("ANTHROPIC_BASE_URL")
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.anthropic.com"
	}
	if cfg.Model == "" {
		cfg.Model = "claude-3-5-sonnet-20241022"
	}
	if cfg.MaxTokens <= 0 {
		cfg.MaxTokens = 1024
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 60 * time.Second
	}

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("APIKey not set (can be set via Config or ANTHROPIC_AUTH_TOKEN environment variable)")
	}

	// Create resty client
	restyClient := resty.New().
		SetBaseURL(cfg.BaseURL).
		SetHeader("x-api-key", cfg.APIKey).
		SetHeader("anthropic-version", "2023-06-01").
		SetHeader("Content-Type", "application/json").
		SetTimeout(cfg.Timeout)

	return &Client{
		config: cfg,
		client: restyClient,
	}, nil
}

// CreateMessage sends multi-turn messages (advanced usage)
func (c *Client) CreateMessage(messages []Message, opts ...Option) (*apiResponse, error) {
	cfg := c.config

	// Apply optional parameter overrides
	for _, opt := range opts {
		opt(&cfg)
	}

	payload := requestPayload{
		Model:     cfg.Model,
		MaxTokens: cfg.MaxTokens,
		Messages:  messages,
	}

	var resp apiResponse
	_, err := c.client.R().
		SetBody(payload).
		SetResult(&resp).
		Post("/v1/messages")

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return &resp, nil
}

// SimplePrompt single-turn conversation shortcut method
func (c *Client) SimplePrompt(prompt string, opts ...Option) (string, error) {
	messages := []Message{{Role: "user", Content: prompt}}
	resp, err := c.CreateMessage(messages, opts...)
	if err != nil {
		return "", err
	}

	for _, block := range resp.Content {
		if block.Type == "text" {
			return block.Text, nil
		}
	}

	return "", fmt.Errorf("no text content in response")
}

// Option functional option pattern for dynamically overriding configuration
type Option func(*Config)

func WithModel(model string) Option {
	return func(c *Config) {
		c.Model = model
	}
}

func WithMaxTokens(maxTokens int) Option {
	return func(c *Config) {
		c.MaxTokens = maxTokens
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}
