package thread

import (
	"github.com/umk/llmservices/internal/service/handlers/client"
	"github.com/umk/llmservices/pkg/client/thread"
)

func GetClient(clientID string) (*thread.Client, error) {
	cl, err := client.GetClient(clientID)
	if err != nil {
		return nil, err
	}

	return (*thread.Client)(cl), nil
}
