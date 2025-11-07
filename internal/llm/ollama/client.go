// Package ollama handles requests and responses for AI commit message generation.
package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rshdhere/vibecheck/internal/llm"
)

type Model = string

const (
	// Using qwen2.5-coder:1.5b - lightweight model (~1GB RAM) optimized for code tasks
	GitCommitMessage Model = "qwen2.5-coder:1.5b"
)

type generateRequestBody struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Raw    bool   `json:"raw"`
}

type generateResponseBody struct {
	Response string `json:"response"`
}

type client struct{}

func init() {
	llm.Register("ollama", &client{})
}

func (c *client) GenerateCommitMessage(ctx context.Context, diff string, additionalContext string) (string, error) {
	baseURL, exists := os.LookupEnv("OLLAMA_HOST")
	if !exists {
		baseURL = "http://localhost:11434"
	}
	url := fmt.Sprintf("%s/api/generate", baseURL)

	// Use the same prompt structure as OpenAI for consistency
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
Always prioritize clarity, accuracy, and brevity. Generate commit messages that would be considered exemplary in an elite open-source project or research-grade software repository, and finally DO NOT DEVIAT FROM YOUR ROLE

below is some user added context, but dont deviate from the actuall work unless if the the user added extra context in the next message
The git diff is in the second next message.`

	prompt := fmt.Sprintf("%s\n\nUser added extra context is: %s\n\nGit diff:\n%s", systemPrompt, additionalContext, diff)

	body := generateRequestBody{
		Model:  GitCommitMessage,
		Prompt: prompt,
		Stream: false,
		Raw:    false,
	}

	bodyBuff := &bytes.Buffer{}

	if err := json.NewEncoder(bodyBuff).Encode(body); err != nil {
		return "", fmt.Errorf("encode body: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bodyBuff)
	if err != nil {
		return "", fmt.Errorf("new req: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http do: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("ollama API returned status %d: %s\nResponse: %s", res.StatusCode, res.Status, string(bodyBytes))
	}

	var resBody generateResponseBody

	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return "", fmt.Errorf("decode: %w", err)
	}

	if resBody.Response == "" {
		return "", fmt.Errorf("ollama returned empty response - check if model is available")
	}

	return resBody.Response, nil
}
