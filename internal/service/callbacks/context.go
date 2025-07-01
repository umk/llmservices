package callbacks

import (
	"context"

	"github.com/umk/jsonrpc2"
)

type ContextKey string

const (
	CtxClient ContextKey = "clients"
)

func Context(ctx context.Context) context.Context {
	return context.WithValue(ctx, CtxClient, new(jsonrpc2.Client))
}

func Client(ctx context.Context) *jsonrpc2.Client {
	return ctx.Value(CtxClient).(*jsonrpc2.Client)
}
