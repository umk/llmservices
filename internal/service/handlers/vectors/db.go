package vectors

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/umk/jsonrpc2"
	"github.com/umk/llmservices/internal/config"
	"github.com/umk/llmservices/pkg/vectorsdb"
)

var databases sync.Map

func CreateDatabase(ctx context.Context, c jsonrpc2.RPCContext) (any, error) {
	var req createDatabaseRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	db := vectorsdb.NewDatabase[json.RawMessage](
		req.VectorLength,
		vectorsdb.WithRepackPercent(config.Cur.RepackPercent),
	)
	if _, loaded := databases.LoadOrStore(req.DatabaseID, db); loaded {
		return nil, errDatabaseAlreadyExists
	}

	resp := createDatabaseResponse{}

	return c.GetResponse(resp)
}

func DeleteDatabase(ctx context.Context, c jsonrpc2.RPCContext) (any, error) {
	var req deleteDatabaseRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	databases.Delete(req.DatabaseID)

	resp := deleteDatabaseResponse{}

	return c.GetResponse(resp)
}

func ReadDatabase(ctx context.Context, c jsonrpc2.RPCContext) (any, error) {
	var req readDatabaseRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	f, err := os.Open(req.FilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	db, err := vectorsdb.Unmarshal[json.RawMessage](r)
	if err != nil {
		return nil, err
	}

	databases.Store(req.DatabaseID, db)

	resp := readDatabaseResponse{}

	return c.GetResponse(resp)
}

func WriteDatabase(ctx context.Context, c jsonrpc2.RPCContext) (any, error) {
	var req writeDatabaseRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	db := getDatabase(req.DatabaseID)
	if db == nil {
		return nil, errDatabaseNotFound
	}

	f, err := os.Create(req.FilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	w := gzip.NewWriter(f)
	defer w.Close()

	if err := vectorsdb.Marshal(w, db); err != nil {
		return nil, err
	}

	resp := writeDatabaseResponse{}

	return c.GetResponse(resp)
}

func getDatabase(databaseID string) *vectorsdb.Database[json.RawMessage] {
	v, ok := databases.Load(databaseID)
	if !ok {
		return nil
	}

	return v.(*vectorsdb.Database[json.RawMessage])
}
