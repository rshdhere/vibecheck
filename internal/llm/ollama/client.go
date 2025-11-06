// Package ollama handles requests and responses for AI commit message generation.
package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Model = string

const (
	GitCommitMessage Model = "tavernari/git-commit-message:sp_commit_mini"
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

func GenerateGitCommit(ctx context.Context, diff string, additionalContext string) (string, error) {
	baseURL, exists := os.LookupEnv("OLLAMA_HOST")
	if !exists {
		baseURL = "http://localhost:11434"
	}
	url := fmt.Sprintf("%s/api/generate", baseURL)
	body := generateRequestBody{
		Model:  GitCommitMessage,
		Prompt: diff,
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

	var resBody generateResponseBody

	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return "", fmt.Errorf("decode: %w", err)
	}

	return resBody.Response, nil
}
