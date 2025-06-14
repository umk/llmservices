package vectors

import (
	"encoding/json"

	"github.com/umk/llmservices/pkg/vectors"
)

/*** Add Vector ***/

type addVectorRequest struct {
	DatabaseID string `json:"database_id" validate:"required"`
	Record     struct {
		ID     vectors.ID      `json:"id" validate:"required"`
		Vector vectors.Vector  `json:"vector" validate:"require,min=1"`
		Data   json.RawMessage `json:"data"`
	} `json:"record" validate:"required"`
}

type addVectorResponse struct {
	ID vectors.ID `json:"id" validate:"required"`
}

/*** Delete Vector ***/

type deleteVectorRequest struct {
	DatabaseID string     `json:"database_id" validate:"required"`
	RecordID   vectors.ID `json:"record_id" validate:"required"`
}

type deleteVectorResponse struct{}

/*** Search Vectors ***/

type searchVectorsRequest struct {
	DatabaseID string           `json:"database_id" validate:"required"`
	Vectors    []vectors.Vector `json:"vectors" validate:"required,min=1,dive,min=1"`
	K          int              `json:"k" validate:"required,min=1"`
}

type searchVectorsResponse struct {
	Records []searchVectorRecord `json:"records"`
}

type searchVectorRecord struct {
	ID   vectors.ID      `json:"id" validate:"required"`
	Data json.RawMessage `json:"data"`
}

/*** Add Vectors Batch ***/

type addVectorsBatchRequest struct {
	DatabaseID string `json:"database_id" validate:"required"`
	Records    []struct {
		ID     vectors.ID      `json:"id" validate:"required"`
		Vector vectors.Vector  `json:"vector" validate:"required,min=1"`
		Data   json.RawMessage `json:"data"`
	} `json:"records" validate:"required,min=1,dive"`
}

type addVectorsBatchResponse struct {
	Records []addVectorsBatchRecord `json:"records"`
}

type addVectorsBatchRecord struct {
	ID vectors.ID `json:"id" validate:"required"`
}

/*** Delete Vectors Batch ***/

type deleteVectorsBatchRequest struct {
	DatabaseID string       `json:"database_id" validate:"required"`
	RecordIDs  []vectors.ID `json:"record_ids" validate:"required,min=1,dive,required"`
}

type deleteVectorsBatchResponse struct{}

/*** Get Similarity ***/

type getSimilarityRequest struct {
	Vector1 vectors.Vector `json:"vector1" validate:"required"`
	Vector2 vectors.Vector `json:"vector2" validate:"required"`
}

type getSimilarityResponse struct {
	Similarity float32 `json:"similarity"`
}
