# Zhipu AI Usage

## Installation

```bash
go get github.com/afeiship/go-claude
```

## Quick Start

### 1. Set Environment Variables

```bash
export ANTHROPIC_AUTH_TOKEN="your-zhipu-api-key"
export ANTHROPIC_BASE_URL="https://open.bigmodel.cn"
```

### 2. Basic Usage

```go
package main

import (
    "fmt"
    "log"
    "github.com/afeiship/go-claude"
)

func main() {
    // Create Zhipu client
    client, err := claude.NewZhipuClient()
    if err != nil {
        log.Fatal(err)
    }

    // Single turn conversation
    response, err := client.SimplePrompt("Hello")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(response)
}
```

### 3. Multi-turn Conversation

```go
messages := []claude.Message{
    {Role: "user", Content: "What is artificial intelligence?"},
    {Role: "assistant", Content: "Artificial intelligence is..."},
    {Role: "user", Content: "Can you explain in more detail?"},
}

resp, err := client.CreateMessage(messages)
if err != nil {
    log.Fatal(err)
}

fmt.Println(resp.Content[0].Text)
```

### 4. Custom Configuration

```go
cfg := claude.WithZhipuDefaults(claude.Config{
    Model:     "glm-4",          // Override default model
    MaxTokens: 2048,            // Custom max tokens
})

client, err := claude.NewClient(cfg)
```

**Default Configuration:**
- Model: `glm-4.5-flash`
- Max tokens: 1024
- Timeout: 60s