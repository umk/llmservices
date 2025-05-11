package main

import (
	"context"
	"log"
	"os"

	"github.com/umk/llmservices/internal/config"
	"github.com/umk/llmservices/internal/jsonrpc"
	"github.com/umk/llmservices/internal/service"
	"github.com/umk/llmservices/internal/validator"
	"github.com/umk/llmservices/pkg/adapter"
	"github.com/umk/llmservices/pkg/vectors"
)

func main() {
	adapter.InitValidator(validator.V)

	if err := config.Init(); err != nil {
		log.Fatalln("Init error:", err)
	}

	vectors.Init(config.C.VectorSize)

	handler := service.Handler()
	server := jsonrpc.NewServer(handler)

	ctx := context.Background()
	if err := server.Run(ctx, os.Stdin, os.Stdout); err != nil {
		log.Fatalln("Error running server:", err)
	}
}
