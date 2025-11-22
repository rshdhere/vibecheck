// Package perplexity handles requests to the Perplexity chat completions API.
package perplexity

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rshdhere/vibecheck/internal/keys"
	"github.com/rshdhere/vibecheck/internal/llm"
)

type client struct{}

func init() {
	llm.Register("perplexity", &client{})
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

type chatResponse struct {
	Choices []struct {
		Message message `json:"message"`
	} `json:"choices"`
}

func (c *client) GenerateCommitMessage(ctx context.Context, diff string, additionalContext string) (string, error) {
	key, exists := keys.GetAPIKey("perplexity")
	if !exists {
		return "", fmt.Errorf("PERPLEXITY_API_KEY environment variable not set")
	}

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

	userPrompt := fmt.Sprintf(
		"Additional context from user:\n%s\n\nGit diff:\n%s",
		additionalContext,
		diff,
	)

	reqBody := chatRequest{
		Model: "sonar",
		Messages: []message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		MaxTokens:   512,
		Temperature: 0.2,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.perplexity.ai/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", key))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("perplexity API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var chatResp chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response choices from Perplexity")
	}

	return chatResp.Choices[0].Message.Content, nil
}
