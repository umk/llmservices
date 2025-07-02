package main

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/umk/jsonrpc2"
	"github.com/umk/llmservices/internal/config"
	"github.com/umk/llmservices/internal/service"
	"github.com/umk/llmservices/internal/service/callbacks"
	"github.com/umk/llmservices/internal/service/handlers"
)

type Runner struct{}

func (r Runner) Run(ctx context.Context, in io.Reader, out io.Writer) error {
	h := service.Handler()

	ctx = handlers.Context(ctx)
	ctx = callbacks.Context(ctx)

	return jsonrpc2.NewHost(in, out, jsonrpc2.WithServer(h)).Run(ctx)
}

func Serve(ctx context.Context) error {
	s := jsonrpc2.NewServer(Runner{})

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(ch)

	go func() {
		<-ch
		s.Close()
	}()

	if config.Cur.Socket != "" {
		return s.ServeFromNetwork(ctx, "unix", config.Cur.Socket)
	} else {
		return s.ServeFromIO(ctx, os.Stdin, os.Stdout)
	}
}
