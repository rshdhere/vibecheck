// Package openai is responsible for all openai-calls
package openai

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
	llm.Register("openai", &client{})
}

func (c *client) GenerateCommitMessage(ctx context.Context, diff string, additionalContext string) (string, error) {
	key, exists := keys.GetAPIKey("openai")
	if !exists {
		return "", fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}
	client := openaisdk.NewClient(
		option.WithAPIKey(key),
	)
	chatCompletion, err := client.Chat.Completions.New(ctx, openaisdk.ChatCompletionNewParams{
		Messages: []openaisdk.ChatCompletionMessageParamUnion{
			openaisdk.SystemMessage(llm.GetSystemPrompt("openai")),

			openaisdk.UserMessage(fmt.Sprintf("User added extra context is: %s", additionalContext)),
			openaisdk.UserMessage(diff),
		},
		Model: openaisdk.ChatModel("gpt-5-mini"),
	})
	if err != nil {
		return fmt.Sprintf("error while prompting to open-ai at: %v", err), err
	}
	return chatCompletion.Choices[0].Message.Content, nil
}
