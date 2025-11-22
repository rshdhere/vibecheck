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
			openaisdk.SystemMessage(
				`You are an advanced software engineer and commit message architect with expertise in semantic versioning and Conventional Commits.
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
Always prioritize clarity, accuracy, and brevity. Generate commit messages that would be considered exemplary in an elite open-source project or research-grade software repository, and finally DO NOT DEVIAT FROM YOUR ROLE

below is some user added context, but dont deviate from the actuall work unless if the the user added extra context in the next message
The git diff is in the second next message.`),

			openaisdk.UserMessage(fmt.Sprintf("User added extra context is: %s", additionalContext)),
			openaisdk.UserMessage(diff),
		},
		Model: openaisdk.ChatModelGPT4oMini,
	})
	if err != nil {
		return fmt.Sprintf("error while prompting to open-ai at: %v", err), err
	}
	return chatCompletion.Choices[0].Message.Content, nil
}
