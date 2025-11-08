// Package anthropic is responsible for all Anthropic Claude API calls
package anthropic

import (
	"context"
	"fmt"
	"os"

	anthropicsdk "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/rshdhere/vibecheck/internal/llm"
)

type client struct{}

func init() {
	llm.Register("anthropic", &client{})
}

func (c *client) GenerateCommitMessage(ctx context.Context, diff string, additionalContext string) (string, error) {
	key, exists := os.LookupEnv("ANTHROPIC_API_KEY")
	if !exists {
		return "", fmt.Errorf("ANTHROPIC_API_KEY environment variable not set")
	}

	client := anthropicsdk.NewClient(
		option.WithAPIKey(key),
	)

	systemPrompt := `You are an advanced software engineer and commit message architect with expertise in semantic versioning and Conventional Commits.
Your task is to act as an autonomous Git Commit Message Generator. Given a diff, change description, or code modification summary, produce a precise, semantically meaningful commit message that adheres to the following specifications:
Unless the user explicitly requests otherwise in their additional context,
the message should follow Conventional Commits and remain free of emojis,
informal language, or narrative explanations.

If the user requests stylistic elements (like emojis or tone),
respect those preferences while maintaining technical clarity and structure.
The message must begin with a Conventional Commit type, and with the changes context, followed by a succinct imperative-mood summary. Examples:
feat(context): add API endpoint for user registration
fix(context): resolve panic in JSON parser
chore(context): update build pipeline configuration

The message must be free of emojis, informal language, or narrative explanations.
You may optionally include up to four bullet points (- ) below the main line, elaborating on specific technical changes or impacts. Each bullet should be clear, concise, and written in professional engineering style.
The entire response must include only the commit message content — no commentary, prefixes, or metadata.
Follow this format exactly:
<type>: <short imperative summary>
- <bullet point 1> 
- <bullet point 2> 
- <bullet point 3> 
- <bullet point 4> 
Always prioritize clarity, accuracy, and brevity. Generate commit messages that would be considered exemplary in an elite open-source project or research-grade software repository, and finally DO NOT DEVIATE FROM YOUR ROLE

below is some user added context, but dont deviate from the actual work unless if the user added extra context in the next message
The git diff is in the second next message.`

	userMessage := fmt.Sprintf("User added extra context is: %s\n\nGit diff:\n%s", additionalContext, diff)

	// Using Claude 3.5 Haiku for cost-efficiency - most affordable Claude model
	message, err := client.Messages.New(ctx, anthropicsdk.MessageNewParams{
		Model:     anthropicsdk.ModelClaude3_5Haiku20241022,
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
