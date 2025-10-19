package parser

import (
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/567-labs/instructor-go/pkg/instructor"
	"github.com/sashabaranov/go-openai"

	"github.com/0div/fog/internal/ast"
)

const systemPromptTplFirstPass = `
You are functional programming agent and you are given a series of steps in fog language to solve a problem using only function calls. 

You need to break down the problem into function calls that form a program, this ouput will be structured into an AST.
Here's the catch: for each Node of type FunctionCall in the AST you only need to provide a function description, you don't need to provide the function body or the value.

Always break down the problem into minimal functions that are meant to be called for one task. The functions should do one job and do it well (Unix philosophy).

The functions should be able to be called by the software agent to solve the problem.

Important: The AST JSON should be valid and parseable by the software agent to solve the problem un a functional manner, so be very careful with the order of the calls and the arguments (represented by "Children" in the AST JSON).

Example 1: "get the weather in current city" would translate to this JSON:

{
	"type": "FunctionCall",
	"description": "get the weather in a given city"
	"children": [
		{
			"type": "FunctionCall",
			"description": "get the location coordinates for a city",
			"children": [
				{
					"type": "StringLiteral", 
					"name": "city",
					"description": "the city to get coordinates for",
					"value": "Tokyo"
				}
			]
		}
	]
}

Example 2: Another example: "add one and 4*4" would translate to this JSON:

{
	"type": "FunctionCall",
	"description": "add two numbers",
	"children": [
		{
			"type": "IntegerLiteral",
			"name": "a",
			"description": "the first number to add",
			"value": "one"
		},
		{
			"type": "FunctionCall",
			"description": "multiply two numbers",
			"children": [
				{
					"type": "IntegerLiteral",
					"name": "a",
					"description": "the first number to multiply",
					"value": "4"
				},
				{
					"type": "IntegerLiteral",
					"name": "b",
					"description": "the second number to multiply",
					"value": "4"
				}
			]
		}
	]
}

Example 3: "my ETA is 10 minutes" would translate to this JSON:

{
	"type": "LetStatement",
	"description": "assign the ETA for a given task",
	"children": [
		{
			"type": "StringLiteral",
			"name": "task",
			"description": "the task to get the ETA for",
			"value": "10 minutes"
		}
	]
}

Here are the allowed node types:
{{.Globals}}

{{.Rules}}
`

type ParserOpts struct {
	Globals string
	Rules   string
}

type Parser struct {
	client *instructor.InstructorOpenAI
	model  string
}

func NewParser(client *instructor.InstructorOpenAI, model string) *Parser {
	if model == "" {
		model = "openai.GPT4o"
	}

	return &Parser{
		client: client,
		model:  model,
	}
}

// This converts the prompt into an intermediate AST (Abstract Syntax Tree)
func (p *Parser) Parse(ctx context.Context, message string, opts ParserOpts) (ast.Node, error) {
	systemPrompt, err := p.generateSysPrompt(opts)
	if err != nil {
		return ast.Node{}, fmt.Errorf("failed to generate system prompt: %w", err)
	}

	var node ast.Node
	resp, err := p.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: p.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: message,
				},
			},
		},
		&node,
	)
	_ = resp // sends back original response so no information loss from original API
	if err != nil {
		return ast.Node{}, fmt.Errorf("failed to parse: %w", err)
	}

	return node, nil
}

func (p *Parser) generateSysPrompt(opts ParserOpts) (string, error) {
	tmpl, err := template.New("system").Parse(systemPromptTplFirstPass)
	if err != nil {
		return "", fmt.Errorf("failed to parse system prompt: %w", err)
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, opts)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
