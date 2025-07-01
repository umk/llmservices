package agent

import (
	"context"

	"github.com/umk/llmservices/internal/service/handlers"
	"github.com/umk/llmservices/pkg/client/agent"
)

func GetClient(ctx context.Context, clientID string) (*agent.Client, error) {
	cl, err := handlers.GetClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	return (*agent.Client)(cl), nil
}
