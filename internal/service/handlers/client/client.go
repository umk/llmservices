package client

import (
	"context"
	"sync"

	"github.com/umk/jsonrpc2"
	"github.com/umk/llmservices/pkg/client"
)

var clients sync.Map

func SetClient(ctx context.Context, c jsonrpc2.RPCContext) (any, error) {
	var req setClientRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	conf, err := getClientConfig(&req.Config)
	if err != nil {
		return nil, err
	}

	cl := client.New(conf)
	clients.Store(req.ClientId, cl)

	var resp setClientResponse

	return c.GetResponse(resp)
}

func getClient(clientId string) *client.Client {
	v, ok := clients.Load(clientId)
	if !ok {
		return nil
	}

	return v.(*client.Client)
}
