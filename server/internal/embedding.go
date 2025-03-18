package internal

import "github.com/umk/phishell/util/marshalx"

type CreateEmbeddingRequest = marshalx.Dynamic[struct {
	// ID of the model to use.
	Model string `json:"model"`
}]

type CreateEmbeddingResponse = marshalx.Dynamic[struct{}]
