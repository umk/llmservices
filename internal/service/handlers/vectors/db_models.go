package vectors

/*** Create Database ***/

type createDatabaseRequest struct {
	DatabaseID   string `json:"database_id" validate:"required"`
	VectorLength int    `json:"vector_length" validate:"required,min=1"`
}

type createDatabaseResponse struct{}

/*** Delete Database ***/

type deleteDatabaseRequest struct {
	DatabaseID string `json:"database_id" validate:"required"`
}

type deleteDatabaseResponse struct{}

/*** Read Database ***/

type readDatabaseRequest struct {
	DatabaseID string `json:"database_id" validate:"required"`
	FilePath   string `json:"file_path" validate:"required"`
}

type readDatabaseResponse struct{}

/*** Write Database ***/

type writeDatabaseRequest struct {
	DatabaseID string `json:"database_id" validate:"required"`
	FilePath   string `json:"file_path" validate:"required"`
}

type writeDatabaseResponse struct{}
