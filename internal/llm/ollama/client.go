// Package ollama handles requests and responses for AI commit message generation.
package ollama

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

type Model = string

const (
	// using gpt-oss:20b for optimal performance and accuracy as it is an free/opensource llm directly from openai
	GitCommitMessage Model = "gpt-oss:20b"
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
	baseURL, exists := keys.GetAPIKey("ollama")
	if !exists {
		baseURL = "http://localhost:11434"
	}
	url := fmt.Sprintf("%s/api/generate", baseURL)

	systemPrompt := llm.GetSystemPrompt("ollama")

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
