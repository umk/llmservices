package handlers

import (
	"context"
	"sync"
)

type ContextKey string

const (
	CtxClients ContextKey = "clients"
)

func Context(ctx context.Context) context.Context {
	return context.WithValue(ctx, CtxClients, new(sync.Map))
}

func Clients(ctx context.Context) *sync.Map {
	return ctx.Value(CtxClients).(*sync.Map)
}
