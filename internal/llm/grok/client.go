// Package grok is responsible for all xAI Grok API calls
package grok

import (
	"context"
	"fmt"

	openaisdk "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/rshdhere/vibecheck/internal/keys"
	"github.com/rshdhere/vibecheck/internal/llm"
)

type client struct{}

func init() {
	llm.Register("grok", &client{})
}

func (c *client) GenerateCommitMessage(ctx context.Context, diff string, additionalContext string) (string, error) {
	key, exists := keys.GetAPIKey("grok")
	if !exists {
		return "", fmt.Errorf("XAI_API_KEY environment variable not set")
	}

	// Grok uses OpenAI-compatible API
	client := openaisdk.NewClient(
		option.WithAPIKey(key),
		option.WithBaseURL("https://api.x.ai/v1"),
	)

	// Using grok-4-1-fast-reasoning model
	chatCompletion, err := client.Chat.Completions.New(ctx, openaisdk.ChatCompletionNewParams{
		Messages: []openaisdk.ChatCompletionMessageParamUnion{
			openaisdk.SystemMessage(llm.GetSystemPrompt("grok")),

			openaisdk.UserMessage(fmt.Sprintf("User added extra context is: %s", additionalContext)),
			openaisdk.UserMessage(diff),
		},
		Model: "grok-4-1-fast-reasoning",
	})
	if err != nil {
		return "", fmt.Errorf("error while prompting to Grok: %w", err)
	}
	return chatCompletion.Choices[0].Message.Content, nil
}
