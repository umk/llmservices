package client

import (
	"context"
	"sync"

	"github.com/umk/jsonrpc2"
	"github.com/umk/llmservices/pkg/client"
)

var clients sync.Map

func SetClientRPC(ctx context.Context, c jsonrpc2.RPCContext) (any, error) {
	var req setClientRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	cl, err := client.New(&req.Config)
	if err != nil {
		return nil, newConfigError(err)
	}

	SetClient(req.ClientID, cl)

	var resp setClientResponse

	return c.GetResponse(resp)
}

func GetClient(clientID string) (*client.Client, error) {
	v, ok := clients.Load(clientID)
	if !ok {
		return nil, errClientNotFound
	}

	return v.(*client.Client), nil
}

func SetClient(clientID string, client *client.Client) {
	clients.Store(clientID, client)
}
