package seed

import (
	"context"
	_ "embed"
	"encoding/json"
	"log/slog"

	"github.com/0div/fog/internal/cfg"
	"github.com/0div/fog/internal/db/postgres"
	"github.com/0div/fog/internal/discovery"
	"github.com/sashabaranov/go-openai"
)

//go:embed seed_data.json
var rawJson []byte

func Seed() {
	slog.Info("Seeding database")

	var functions []discovery.Function
	if err := json.Unmarshal(rawJson, &functions); err != nil {
		slog.Error("error unmarshalling seed_data.json", "err", err)
		panic(err)
	}
	slog.Info("Loaded seed data", "count", len(functions))

	apiKey := cfg.MustStr("OPENAI_API_KEY")
	openaiClient := openai.NewClient(apiKey)

	pg := postgres.NewPostgresDB()

	embedder := discovery.NewEmbeddings(openaiClient, pg)

	for _, function := range functions {
		err := embedder.EmbedFunction(context.Background(), function)
		if err != nil {
			slog.Error("error embedding function", "err", err)
			continue
		}
		slog.Debug("embedded function", "function", function.Name)
	}
	slog.Info("Database seeded successfully", "row_count", len(functions))
}
