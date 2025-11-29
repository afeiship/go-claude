package claude

import (
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

// Config 定义客户端配置
type Config struct {
	APIKey    string        // 必填
	BaseURL   string        // 可选，默认 "https://api.anthropic.com"
	Model     string        // 可选，默认 "claude-3-5-sonnet-20241022"
	MaxTokens int           // 可选，默认 1024
	Timeout   time.Duration // 可选，默认 60s
}

// Client 是 Claude API 客户端
type Client struct {
	config Config
	client *resty.Client
}

// Message 表示一条对话消息
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

// NewClient 从 Config 创建客户端
func NewClient(cfg Config) (*Client, error) {
	// 使用环境变量作为 fallback（可选，你也可以完全禁用）
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
		return nil, fmt.Errorf("APIKey 未设置（可通过 Config 或 ANTHROPIC_AUTH_TOKEN 环境变量）")
	}

	// 创建 resty 客户端
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

// CreateMessage 发送多轮消息（高级用法）
func (c *Client) CreateMessage(messages []Message, opts ...Option) (*apiResponse, error) {
	cfg := c.config

	// 应用可选参数覆盖
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
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	return &resp, nil
}

// SimplePrompt 单轮对话快捷方法
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

	return "", fmt.Errorf("响应中无文本内容")
}

// Option 函数式选项模式，用于动态覆盖配置
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
