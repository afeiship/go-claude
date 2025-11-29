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
	"your-project/claude"
)

func main() {
	// 方式 1：完全通过配置传入（推荐用于生产）
	cfg := claude.Config{
		APIKey:    "sk-ant-xxxxxxxx",           // 来自你的配置中心/secret 管理
		BaseURL:   "https://api.anthropic.com", // 可省略，默认值
		Model:     "claude-3-haiku-20240307",   // 可动态指定
		MaxTokens: 512,
		Timeout:   30 * time.Second,
	}

	client, err := claude.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// 使用默认配置
	reply, err := client.SimplePrompt("你好，Claude！")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("回复1:", reply)

	// 临时覆盖模型和 max_tokens（使用 Option）
	reply2, err := client.SimplePrompt(
		"用一句话解释量子计算",
		claude.WithModel("claude-3-5-sonnet-20241022"),
		claude.WithMaxTokens(256),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("回复2:", reply2)
}
```