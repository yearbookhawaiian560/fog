package repl

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/567-labs/instructor-go/pkg/instructor"
	"github.com/sashabaranov/go-openai"

	"github.com/0div/fog/internal/ast"
	"github.com/0div/fog/internal/cfg"
	"github.com/0div/fog/internal/db/postgres"
	"github.com/0div/fog/internal/discovery"
	"github.com/0div/fog/internal/interpreter"
	"github.com/0div/fog/internal/parser"
)

type ReplOpts struct {
	OpenAIClient *openai.Client
	Model        string
	Postgres     *postgres.PostgresDB
}

type Repl struct {
	embedder    *discovery.Embeddings
	parser      *parser.Parser
	interpreter *interpreter.Interpreter
}

func NewRepl(opts ReplOpts) *Repl {
	slog.Debug("creating repl", "opts", opts)
	embedder := discovery.NewEmbeddings(opts.OpenAIClient, opts.Postgres)

	parserInstructorClient := instructor.FromOpenAI(
		opts.OpenAIClient,
		instructor.WithMode(instructor.ModeJSON),
		instructor.WithMaxRetries(3),
	)

	parser_ := parser.NewParser(parserInstructorClient, opts.Model)

	interpreter, err := interpreter.NewInterpreter()
	if err != nil {
		slog.Error("error creating interpreter", "err", err)
		os.Exit(1)
	}

	return &Repl{
		embedder:    embedder,
		parser:      parser_,
		interpreter: interpreter,
	}
}

const introMessage = `
	 ░░▒▒▓▓██▓▓▒▒░   
	░▒▓█████████▓▒░   
	▒▓▓▓ F O G ▓▓▓▒  
	░▒▓█████████▓▒░   
	 ░░▒▒▓▓██▓▓▒▒░   
`

func (m *Repl) Run(ctx context.Context) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println()
	fmt.Println(introMessage)
	fmt.Println()

	// in development mode, this first input is used to test the REPL without having to enter it manually
	firstInput := "get weather in my current city"
	firstInputDone := !cfg.Bool("development") // only enter the first input once in development mode

	for {
		fmt.Print("> ")
		if !firstInputDone {
			firstInputDone = true
			fmt.Printf(" %s\n", firstInput)
			// On very first loop, inject "get weather"
			scanner = bufio.NewScanner(strings.NewReader(firstInput + "\n"))
		} else {
			scanner = bufio.NewScanner(os.Stdin)
		}

		if !scanner.Scan() {
			break
		}
		prompt := scanner.Text()

		fmt.Println()
		fmt.Println("╔════════════════════════════════════════════════════════════╗")
		fmt.Println("║         Parsing prompt into the Intermediate AST...        ║")
		fmt.Println("╚════════════════════════════════════════════════════════════╝")
		fmt.Println()

		node, err := m.parser.Parse(ctx, prompt, parser.ParserOpts{
			Globals: ast.NodeTypesAsString(),
		})
		if err != nil {
			slog.Error("parser error", "error", err)
			continue
		}
		slog.Debug("parser completed successfully")

		json1, err := json.MarshalIndent(node, "", "  ")
		if err != nil {
			slog.Error("json error", "error", err)
			continue
		}
		astJSON1 := string(json1)
		fmt.Printf("%s\n", astJSON1)

		fmt.Println()
		fmt.Println("╔════════════════════════════════════════════════════════════╗")
		fmt.Println("║         Compiling Intermediate AST into FOG AST...         ║")
		fmt.Println("╚════════════════════════════════════════════════════════════╝")
		fmt.Println()

		traverseNode(m.embedder, &node, 0, cfg.Bool("debug"))

		json2, err := json.MarshalIndent(node, "", "  ")
		if err != nil {
			slog.Error("json error", "error", err)
			continue
		}
		astJSON2 := string(json2)
		fmt.Printf("%s\n", astJSON2)

		fmt.Println()
		fmt.Println("╔════════════════════════════════════════════════════════════╗")
		fmt.Println("║           Evaluating the Intermediate AST...               ║")
		fmt.Println("╚════════════════════════════════════════════════════════════╝")
		fmt.Println()

		result, err := m.interpreter.Eval(ctx, node)
		if err != nil {
			slog.Error("interpreter error", "error", err)
			continue
		}
		fmt.Printf("Result:\t%+v\n", result.Export())
		fmt.Println()
		fmt.Println()
	}

	if err := scanner.Err(); err != nil {
		slog.Error("scanner error", "error", err)
	}

	fmt.Println("Exiting parser mode...")
	fmt.Println()
}

func traverseNode(embedder *discovery.Embeddings, node *ast.Node, depth int, debug bool) {
	tabs := ""
	for range depth {
		tabs += "\t"
	}

	if node.Description != "" && node.Type == "FunctionCall" {
		slog.Debug("processing function call", "description", node.Description, "type", node.Type)
		cosineSimilarity, err := embedder.CosineSimilarity(context.Background(), node.Description)
		if err != nil {
			slog.Error("error getting cosine similarity", "error", err)
			os.Exit(1)
		}
		distance := "N/A"
		if len(cosineSimilarity) > 0 {
			distance = fmt.Sprintf("%f", cosineSimilarity[0].Distance)
			node.Distance = cosineSimilarity[0].Distance.(float64)
			node.JS = cosineSimilarity[0].Js
			slog.Debug("cosine similarity result", "distance", distance, "functionName", cosineSimilarity[0].FunctionName)
		}

		functionName := "N/A"
		if len(cosineSimilarity) > 0 {
			functionName = cosineSimilarity[0].FunctionName
		}
		node.Name = functionName

		if debug {
			fmt.Printf("%s<%s> Description: %s [`%s` cosine similarity == %v]\n", tabs, node.Type, node.Description, node.Name, distance)
			fmt.Printf("%s<JS Code>: %s\n", tabs, node.JS)
		}
	}

	if strings.HasSuffix(string(node.Type), "Literal") {
		slog.Debug("processing literal", "type", node.Type, "name", node.Name, "value", node.Value)
		if debug {
			fmt.Printf("%s<ARG %s> %s = %s\n", tabs, node.Type, node.Name, node.Value)
		}
	}

	// Recursively traverse children (left to right)
	for _, child := range node.Children {
		traverseNode(embedder, child, depth+1, debug)
	}
}
