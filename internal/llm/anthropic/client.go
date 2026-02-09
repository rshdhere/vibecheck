// Package anthropic is responsible for all Anthropic Claude API calls
package anthropic

import (
	"context"
	"fmt"

	anthropicsdk "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/rshdhere/vibecheck/internal/keys"
	"github.com/rshdhere/vibecheck/internal/llm"
)

type client struct{}

func init() {
	llm.Register("anthropic", &client{})
}

func (c *client) GenerateCommitMessage(ctx context.Context, diff string, additionalContext string) (string, error) {
	key, exists := keys.GetAPIKey("anthropic")
	if !exists {
		return "", fmt.Errorf("ANTHROPIC_API_KEY environment variable not set")
	}

	client := anthropicsdk.NewClient(
		option.WithAPIKey(key),
	)

	systemPrompt := llm.GetSystemPrompt("anthropic")

	userMessage := fmt.Sprintf("User added extra context is: %s\n\nGit diff:\n%s", additionalContext, diff)

	// Using Anthropic's current Haiku alias for low-cost commit-message generation.
	message, err := client.Messages.New(ctx, anthropicsdk.MessageNewParams{
		Model:     anthropicsdk.Model("claude-haiku-4-5-20251001"),
		MaxTokens: 1024,
		Messages: []anthropicsdk.MessageParam{
			anthropicsdk.NewUserMessage(anthropicsdk.NewTextBlock(userMessage)),
		},
		System: []anthropicsdk.TextBlockParam{
			{
				Text: systemPrompt,
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("error while prompting to Anthropic: %w", err)
	}

	if len(message.Content) == 0 {
		return "", fmt.Errorf("no response generated from Anthropic")
	}

	return message.Content[0].Text, nil
}
