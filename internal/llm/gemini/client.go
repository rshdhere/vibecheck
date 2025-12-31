// Package gemini is responsible for all Google Gemini API calls
package gemini

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"github.com/rshdhere/vibecheck/internal/keys"
	"github.com/rshdhere/vibecheck/internal/llm"
	"google.golang.org/api/option"
)

type client struct{}

func init() {
	llm.Register("gemini", &client{})
}

func (c *client) GenerateCommitMessage(ctx context.Context, diff string, additionalContext string) (string, error) {
	key, exists := keys.GetAPIKey("gemini")
	if !exists {
		return "", fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(key))
	if err != nil {
		return "", fmt.Errorf("create gemini client: %w", err)
	}
	defer client.Close()

	// Using gemini-2.5-flash for best performance and cost-efficiency
	model := client.GenerativeModel("gemini-2.5-flash")

	// Configure model parameters for better responses
	model.SetTemperature(0.7)
	model.SetTopK(40)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(1024)

	// Relax safety settings for commit messages (they're just code diffs)
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockOnlyHigh,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockOnlyHigh,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockOnlyHigh,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockOnlyHigh,
		},
	}

	// Gemini works better with system instructions set on the model
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text(`You are a commit message generator. Analyze the git diff and generate a conventional commit message.
Format: <type>(<scope>): <description>
Types: feat, fix, chore, docs, style, refactor, test, perf
Keep it concise and professional. Add 2-4 bullet points for details.`),
		},
	}

	prompt := fmt.Sprintf("User context: %s\n\nGenerate a conventional commit message for this git diff:\n\n%s", additionalContext, diff)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("generate content: %w", err)
	}

	// Check if response was blocked by safety filters
	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("gemini returned no candidates (possibly blocked by safety filters)")
	}

	candidate := resp.Candidates[0]

	// Check for content filtering
	if candidate.FinishReason != genai.FinishReasonStop && candidate.FinishReason != genai.FinishReasonMaxTokens {
		return "", fmt.Errorf("gemini response blocked: finish reason = %v", candidate.FinishReason)
	}

	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return "", fmt.Errorf("gemini returned empty content")
	}

	return fmt.Sprintf("%v", candidate.Content.Parts[0]), nil
}
