package discovery

import (
	"context"
	"errors"

	"github.com/0div/fog/internal/db/postgres"
	"github.com/pgvector/pgvector-go"
	"github.com/sashabaranov/go-openai"
)

type Function struct {
	Name        string
	Description string
	JS          string
}

type Embeddings struct {
	client *openai.Client
	pg     *postgres.PostgresDB
}

func NewEmbeddings(client *openai.Client, pg *postgres.PostgresDB) *Embeddings {
	return &Embeddings{
		client: client,
		pg:     pg,
	}
}

func (e *Embeddings) EmbedText(text string) (pgvector.Vector, error) {
	embeddingResponse, err := e.client.CreateEmbeddings(context.Background(), openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.SmallEmbedding3,
	})
	if err != nil {
		return pgvector.Vector{}, err
	}

	embeddingVector := embeddingResponse.Data[0].Embedding
	if len(embeddingVector) == 0 {
		return pgvector.Vector{}, errors.New("embedding vector is empty")
	}
	return pgvector.NewVector(embeddingVector), nil
}

func (e *Embeddings) EmbedFunction(ctx context.Context, function Function) error {
	embeddingVector, err := e.EmbedText(function.Description)
	if err != nil {
		return err
	}

	_, err = e.pg.Q.CreateFunction(ctx, postgres.CreateFunctionParams{
		FunctionName: function.Name,
		Embedding:    embeddingVector,
		Js:           function.JS,
	})

	return err
}

func (e *Embeddings) CosineSimilarity(ctx context.Context, text string) ([]postgres.CosineSimilarityRow, error) {
	embeddingVector, err := e.EmbedText(text)
	if err != nil {
		return nil, err
	}

	return e.pg.Q.CosineSimilarity(ctx, postgres.CosineSimilarityParams{
		Embedding: embeddingVector,
		Limit:     10,
	})
}
