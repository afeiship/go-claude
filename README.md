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

	"github.com/afeiship/go-claude"
)

func main() {
	client := claude.NewClient("your-api-key-here")

	// Simple prompt - just get a text response
	response, err := client.SimplePrompt("Hello, Claude! How are you?")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}
```