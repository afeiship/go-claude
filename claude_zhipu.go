package claude

import (
	"fmt"
	"os"
	"time"
)

// NewZhipuClient creates a zhipu client from environment variables (80% use case)
func NewZhipuClient() (*Client, error) {
	cfg := Config{
		Model:     "glm-4.5-flash",
		MaxTokens: 1024,
		Timeout:   60 * time.Second,
	}

	// Use existing environment variables
	cfg.APIKey = os.Getenv("ANTHROPIC_AUTH_TOKEN")
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_AUTH_TOKEN not set")
	}

	cfg.BaseURL = os.Getenv("ANTHROPIC_BASE_URL")
	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("ANTHROPIC_BASE_URL not set")
	}

	return NewClient(cfg)
}

// WithZhipuDefaults creates a config with zhipu-specific defaults
func WithZhipuDefaults(cfg Config) Config {
	if cfg.BaseURL == "" {
		cfg.BaseURL = os.Getenv("ANTHROPIC_BASE_URL")
	}
	if cfg.Model == "" {
		cfg.Model = "glm-4.5-flash"
	}
	if cfg.APIKey == "" {
		cfg.APIKey = os.Getenv("ANTHROPIC_AUTH_TOKEN")
	}
	return cfg
}