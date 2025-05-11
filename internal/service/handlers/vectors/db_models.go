package vectors

/*** Create Database ***/

type createDatabaseRequest struct {
	DatabaseId   string `json:"database_id" validate:"required"`
	VectorLength int    `json:"vector_length" validate:"required,min=1"`
}

type createDatabaseResponse struct{}

/*** Delete Database ***/

type deleteDatabaseRequest struct {
	DatabaseId string `json:"database_id" validate:"required"`
}

type deleteDatabaseResponse struct{}
