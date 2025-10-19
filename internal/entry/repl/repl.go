package repl

import (
	"context"

	"github.com/0div/fog/internal/cfg"
	"github.com/0div/fog/internal/db/postgres"
	"github.com/0div/fog/internal/repl"
	"github.com/sashabaranov/go-openai"
)

func RunFogREPL() {
	apiKey := cfg.MustStr("OPENAI_API_KEY")
	model := cfg.Str("OPENAI_MODEL")
	config := openai.DefaultConfig(apiKey)

	openaiClient := openai.NewClientWithConfig(config)

	pg := postgres.NewPostgresDB()

	repl := repl.NewRepl(repl.ReplOpts{
		OpenAIClient: openaiClient,
		Model:        model,
		Postgres:     pg,
	})

	repl.Run(context.Background())
}
