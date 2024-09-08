package openai

import "errors"

var (
	ErrInvalidTarif = errors.New("invalid tarif")
	ErrNoEmbedding  = errors.New("no embedding")
)
