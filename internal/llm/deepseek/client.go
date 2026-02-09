// Package deepseek is responsible for all DeepSeek API calls
package deepseek

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
	llm.Register("deepseek", &client{})
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type chatResponse struct {
	Choices []struct {
		Message message `json:"message"`
	} `json:"choices"`
}

func (c *client) GenerateCommitMessage(ctx context.Context, diff string, additionalContext string) (string, error) {
	key, exists := keys.GetAPIKey("deepseek")
	if !exists {
		return "", fmt.Errorf("DEEPSEEK_API_KEY environment variable not set")
	}

	systemPrompt := llm.GetSystemPrompt("deepseek")

	reqBody := chatRequest{
		// DeepSeek v3.2 is exposed under the "deepseek-chat" model ID in the chat API.
		Model: "deepseek-chat",
		Messages: []message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: fmt.Sprintf("User added extra context is: %s", additionalContext)},
			{Role: "user", Content: diff},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", key))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("deepseek API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var chatResp chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response choices from DeepSeek")
	}

	return chatResp.Choices[0].Message.Content, nil
}
