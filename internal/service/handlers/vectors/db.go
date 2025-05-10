package vectors

import (
	"context"
	"sync"

	"github.com/umk/llmservices/internal/config"
	"github.com/umk/llmservices/internal/jsonrpc"
	"github.com/umk/llmservices/pkg/vectorsdb"
)

var databases sync.Map

func CreateDatabase(ctx context.Context, c jsonrpc.RPCContext) (any, error) {
	var req createDatabaseRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	db := vectorsdb.NewDatabase[any](
		req.VectorLength,
		vectorsdb.WithRepackPercent(config.C.RepackPercent),
	)
	if _, loaded := databases.LoadOrStore(req.DatabaseId, &db); loaded {
		return nil, errDatabaseAlreadyExists
	}

	resp := createDatabaseResponse{}

	return c.GetResponse(resp)
}

func DeleteDatabase(ctx context.Context, c jsonrpc.RPCContext) (any, error) {
	var req deleteDatabaseRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	databases.Delete(req.DatabaseId)

	resp := deleteDatabaseResponse{}

	return c.GetResponse(resp)
}

func getDatabase(databaseId string) *vectorsdb.Database[any] {
	v, ok := databases.Load(databaseId)
	if !ok {
		return nil
	}

	return v.(*vectorsdb.Database[any])
}
