package parser

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

const systemPromptSteps = `
You are a helpful assistant that distills a prompt into a series of steps for a sofware agent.
These steps should be a series of instructions that will then be used to generate a response from the sofware agent by completing them one by one.
Sometimes the prompt might be ambiguous or incomplete, in which case you can just say "store message" and move on to the next step.

The steps should be descibed in a way that is easy to understand and follow and should be in the following format:

- Step 1: [Step 1 description]
- Step 2: [Step 2 description]
- ...
`

type Steps struct {
	OpenAIClient *openai.Client
	Model        string
}

func NewSteps(openAIClient *openai.Client, model string) *Steps {
	return &Steps{
		OpenAIClient: openAIClient,
		Model:        model,
	}
}

func (s *Steps) Run(ctx context.Context, prompt string) (string, error) {

	resp, err := s.OpenAIClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPromptSteps,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	})

	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}
