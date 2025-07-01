package thread

import (
	"context"

	"github.com/umk/llmservices/internal/service/handlers"
	"github.com/umk/llmservices/pkg/client/thread"
)

func GetClient(ctx context.Context, clientID string) (*thread.Client, error) {
	cl, err := handlers.GetClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	return (*thread.Client)(cl), nil
}
