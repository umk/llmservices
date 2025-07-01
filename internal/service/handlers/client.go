package handlers

import (
	"context"
	"sync"

	"github.com/umk/jsonrpc2"
	"github.com/umk/llmservices/pkg/client"
)

var globalClients sync.Map

func SetClientRPC(ctx context.Context, c jsonrpc2.RPCContext) (any, error) {
	var req SetClientRequest
	if err := c.Request(&req); err != nil {
		return nil, err
	}

	cl, err := client.New(&req.Config)
	if err != nil {
		return nil, newConfigError(err)
	}

	Clients(ctx).Store(req.ClientID, cl)

	var resp SetClientResponse

	return c.Response(resp)
}

func GetClient(ctx context.Context, clientID string) (*client.Client, error) {
	if v, ok := Clients(ctx).Load(clientID); ok {
		return v.(*client.Client), nil
	}

	if v, ok := globalClients.Load(clientID); ok {
		return v.(*client.Client), nil
	}

	return nil, errClientNotFound
}

func SetClient(ctx context.Context, clientID string, client *client.Client) {
	Clients(ctx).Store(clientID, client)
}

func SetGlobalClient(clientID string, client *client.Client) {
	globalClients.Store(clientID, client)
}
