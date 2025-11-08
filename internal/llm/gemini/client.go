// Package gemini is responsible for all Google Gemini API calls
package gemini

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/rshdhere/vibecheck/internal/llm"
	"google.golang.org/api/option"
)

type client struct{}

func init() {
	llm.Register("gemini", &client{})
}

func (c *client) GenerateCommitMessage(ctx context.Context, diff string, additionalContext string) (string, error) {
	key, exists := os.LookupEnv("GEMINI_API_KEY")
	if !exists {
		return "", fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(key))
	if err != nil {
		return "", fmt.Errorf("create gemini client: %w", err)
	}
	defer client.Close()

	// Using gemini-1.5-flash for cost-efficiency and good performance
	model := client.GenerativeModel("gemini-1.5-flash")

	// Configure model parameters for better responses
	model.SetTemperature(0.7)
	model.SetTopK(40)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(1024)

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

	prompt := fmt.Sprintf("%s\n\nUser added extra context is: %s\n\nGit diff:\n%s", systemPrompt, additionalContext, diff)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response generated from Gemini")
	}

	return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]), nil
}
