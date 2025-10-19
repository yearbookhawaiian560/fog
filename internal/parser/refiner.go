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

const systemPromptTplRefine = `
You are functional programming developer and you are given the first version of an AST JSON that needs to be validated and refined. 

The initial prompt was: 

"""
{{.InitialPrompt}}
"""

The initial AST is:

"""
{{.ASTJSON}}
"""
You need to more accurately break down the problem into function calls that form a program, this ouput will be structured into an AST as well.

Things you want to optimize:
- Break down the problem into more functions if a function solves more thant one small task
- Functions that are not needed 
- Functions are too vague
- Redundant function calls
`

type RefinerOpts struct {
	InitialPrompt string
	ASTJSON       string
}

type Refiner struct {
	client *instructor.InstructorOpenAI
	model  string
}

func NewRefiner(client *instructor.InstructorOpenAI, model string) *Refiner {
	if model == "" {
		model = "openai.GPT4o"
	}

	return &Refiner{
		client: client,
		model:  model,
	}
}

func (p *Refiner) Refine(ctx context.Context, opts RefinerOpts) (ast.Node, error) {
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
					Content: opts.InitialPrompt,
				},
			},
		},
		&node,
	)
	_ = resp // sends back original response so no information loss from original API
	if err != nil {
		return ast.Node{}, fmt.Errorf("failed to refine: %w", err)
	}

	return node, nil
}

func (p *Refiner) generateSysPrompt(opts RefinerOpts) (string, error) {
	tmpl, err := template.New("system").Parse(systemPromptTplRefine)
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
