package main

import (
	"context"
	"log"
	"os"

	"github.com/umk/jsonrpc2"
	"github.com/umk/llmservices/internal/config"
	"github.com/umk/llmservices/internal/service"
	"github.com/umk/llmservices/internal/service/client"
	"github.com/umk/llmservices/pkg/adapter"
	"github.com/umk/llmservices/pkg/vectors"
)

func main() {
	adapter.InitValidator(jsonrpc2.Val)

	if err := initConfig(); err != nil {
		log.Fatalln("Init error:", err)
	}

	vectors.Init(config.Cur.VectorBufSize)

	handler := service.Handler()
	host := jsonrpc2.NewHost(
		os.Stdin,
		os.Stdout,
		jsonrpc2.WithClient(&client.Client),
		jsonrpc2.WithServer(handler),
	)

	ctx := context.Background()
	if err := host.Run(ctx); err != nil {
		log.Fatalln("Error running server:", err)
	}
}
