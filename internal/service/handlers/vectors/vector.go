package vectors

import (
	"context"
	"encoding/json"

	"github.com/umk/llmservices/internal/jsonrpc"
	"github.com/umk/llmservices/pkg/vectors"
	"github.com/umk/llmservices/pkg/vectorsdb"
)

func AddVector(ctx context.Context, c jsonrpc.RPCContext) (any, error) {
	var req addVectorRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	db := getDatabase(req.DatabaseId)
	if db == nil {
		return nil, errDatabaseNotFound
	}

	r := db.Add(vectorsdb.Record[json.RawMessage]{
		Vector: req.Record.Vector,
		Data:   req.Record.Data,
	})

	resp := addVectorResponse{
		Id: r.Id,
	}

	return c.GetResponse(resp)
}

func DeleteVector(ctx context.Context, c jsonrpc.RPCContext) (any, error) {
	var req deleteVectorRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	db := getDatabase(req.DatabaseId)
	if db == nil {
		return nil, errDatabaseNotFound
	}

	db.Delete(req.RecordId)

	resp := deleteVectorResponse{}

	return c.GetResponse(resp)
}

func SearchVectors(ctx context.Context, c jsonrpc.RPCContext) (any, error) {
	var req searchVectorsRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	db := getDatabase(req.DatabaseId)
	if db == nil {
		return nil, errDatabaseNotFound
	}

	rs := db.Get(req.Vectors, req.K)

	resp := searchVectorsResponse{
		Records: make([]searchVectorRecord, 0, len(rs)),
	}

	for _, r := range rs {
		resp.Records = append(resp.Records, searchVectorRecord{
			Id:   r.Id,
			Data: r.Data,
		})
	}

	return c.GetResponse(resp)
}

func AddVectorsBatch(ctx context.Context, c jsonrpc.RPCContext) (any, error) {
	var req addVectorsBatchRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	db := getDatabase(req.DatabaseId)
	if db == nil {
		return nil, errDatabaseNotFound
	}

	records := make([]vectorsdb.Record[json.RawMessage], len(req.Records))
	for i, r := range req.Records {
		records[i] = vectorsdb.Record[json.RawMessage]{
			Vector: r.Vector,
			Data:   r.Data,
		}
	}

	rs := db.AddBatch(records)

	resp := addVectorsBatchResponse{
		Records: make([]addVectorsBatchRecord, 0, len(rs)),
	}

	for _, r := range rs {
		resp.Records = append(resp.Records, addVectorsBatchRecord{Id: r.Id})
	}

	return c.GetResponse(resp)
}

func DeleteVectorsBatch(ctx context.Context, c jsonrpc.RPCContext) (any, error) {
	var req deleteVectorsBatchRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	db := getDatabase(req.DatabaseId)
	if db == nil {
		return nil, errDatabaseNotFound
	}

	db.DeleteBatch(req.RecordIds)

	resp := deleteVectorsBatchResponse{}

	return c.GetResponse(resp)
}

func GetSimilarity(ctx context.Context, c jsonrpc.RPCContext) (any, error) {
	var req getSimilarityRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	if len(req.Vector1) != len(req.Vector2) {
		return nil, errVectorsLengthMismatch
	}

	resp := getSimilarityResponse{
		Similarity: vectors.Similarity(req.Vector1, req.Vector2),
	}

	return c.GetResponse(resp)
}
