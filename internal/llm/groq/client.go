// Package groq is responsible for all Groq API calls
package groq

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
	llm.Register("groq", &client{})
}

func (c *client) GenerateCommitMessage(ctx context.Context, diff string, additionalContext string) (string, error) {
	key, exists := keys.GetAPIKey("groq")
	if !exists {
		return "", fmt.Errorf("GROQ_API_KEY environment variable not set")
	}

	// Groq uses OpenAI-compatible API
	client := openaisdk.NewClient(
		option.WithAPIKey(key),
		option.WithBaseURL("https://api.groq.com/openai/v1"),
	)

	// Using openai/gpt-oss-20b for excellent performance and speed
	chatCompletion, err := client.Chat.Completions.New(ctx, openaisdk.ChatCompletionNewParams{
		Messages: []openaisdk.ChatCompletionMessageParamUnion{
			openaisdk.SystemMessage(llm.GetSystemPrompt("groq")),

			openaisdk.UserMessage(fmt.Sprintf("User added extra context is: %s", additionalContext)),
			openaisdk.UserMessage(diff),
		},
		Model: "openai/gpt-oss-20b",
	})
	if err != nil {
		return "", fmt.Errorf("error while prompting to Groq: %w", err)
	}
	return chatCompletion.Choices[0].Message.Content, nil
}
