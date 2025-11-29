# go-claude
> A lightweight, idiomatic Go client for the Anthropic Claude API.

## installation
```sh
go get -u github.com/afeiship/go-claude
```

## usage
```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/afeiship/go-claude"
)

func main() {
	// Method 1: Pass configuration completely (recommended for production)
	cfg := claude.Config{
		APIKey:    "sk-ant-xxxxxxxx",           // From your config center/secret management
		BaseURL:   "https://api.anthropic.com", // Optional, has default value
		Model:     "claude-3-haiku-20240307",   // Can be dynamically specified
		MaxTokens: 512,
		Timeout:   30 * time.Second,
	}

	client, err := claude.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Use default configuration
	reply, err := client.SimplePrompt("Hello, Claude!")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Reply 1:", reply)

	// Temporarily override model and max_tokens (using Option)
	reply2, err := client.SimplePrompt(
		"Explain quantum computing in one sentence",
		claude.WithModel("claude-3-5-sonnet-20241022"),
		claude.WithMaxTokens(256),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Reply 2:", reply2)
}
```