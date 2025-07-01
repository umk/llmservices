package main

import (
	"context"
	"log"

	"github.com/umk/jsonrpc2"
	"github.com/umk/llmservices/internal/config"
	"github.com/umk/llmservices/pkg/adapter"
)

func main() {
	adapter.InitValidator(jsonrpc2.Val)

	if err := config.Init(); err != nil {
		log.Fatalln("Init error:", err)
	}

	if err := Serve(context.Background()); err != nil {
		log.Fatalln("Error running server:", err)
	}
}
