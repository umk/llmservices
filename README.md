# llmservices

[![Go Reference](https://pkg.go.dev/badge/github.com/umk/llmservices.svg)](https://pkg.go.dev/github.com/umk/llmservices)

An in-memory vector store and JSON-RPC service via Stdio that uses a brute-force approach for searching through vectors.

## Service Methods

| Method | Description | Models |
|--------|-------------|:--------:|
| `createDatabase` | Creates a new vector database with specified ID and vector length | [↗](internal/service/handlers/vectors/db_models.go) |
| `deleteDatabase` | Deletes a database by ID | [↗](internal/service/handlers/vectors/db_models.go) |
| `addVector` | Adds a single vector to a database | [↗](internal/service/handlers/vectors/vector_models.go) |
| `deleteVector` | Deletes a vector from a database by ID | [↗](internal/service/handlers/vectors/vector_models.go) |
| `addVectorsBatch` | Adds multiple vectors to a database in a batch operation | [↗](internal/service/handlers/vectors/vector_models.go) |
| `deleteVectorsBatch` | Deletes multiple vectors from a database in a batch operation | [↗](internal/service/handlers/vectors/vector_models.go) |
| `searchVectors` | Searches for vectors in a database that are similar to the provided vectors | [↗](internal/service/handlers/vectors/vector_models.go) |
| `getSimilarity` | Computes the cosine similarity between two vectors | [↗](internal/service/handlers/vectors/vector_models.go) |
| `setClient` | Configures and initializes a client with specified settings | [↗](internal/service/handlers/client/models.go) |
| `getCompletion` | Retrieves AI model completions using the configured client | [↗](internal/service/handlers/client/models.go) [↗](pkg/adapter/completion.go) [↗](pkg/adapter/message.go) [↗](pkg/adapter/tool.go) [↗](pkg/adapter/content.go) |
| `getEmbeddings` | Generates vector embeddings for input text using the configured client | [↗](internal/service/handlers/client/models.go) [↗](pkg/adapter/embeddings.go) |
| `getStatistics` | Returns statistics about the client, including bytes per token | [↗](internal/service/handlers/client/models.go) |
