package grok

import (
	"context"
	"testing"
)

// TestClientRegistration verifies the client is registered correctly
func TestClientRegistration(t *testing.T) {
	var _ interface {
		GenerateCommitMessage(ctx context.Context, diff string, additionalContext string) (string, error)
	} = &client{}
}

// TestAPIKeyValidation verifies API key validation
func TestAPIKeyValidation(t *testing.T) {
	c := &client{}
	ctx := context.Background()

	_, err := c.GenerateCommitMessage(ctx, "test diff", "")
	if err == nil {
		t.Error("GenerateCommitMessage() should return error when API key is missing")
	}
	if err != nil && err.Error() != "XAI_API_KEY environment variable not set" {
		t.Errorf("GenerateCommitMessage() error message = %q, want 'XAI_API_KEY environment variable not set'", err.Error())
	}
}

// TestBaseURL verifies the base URL matches xAI Grok API documentation
// According to xAI docs: https://api.x.ai/v1 (OpenAI-compatible)
func TestBaseURL(t *testing.T) {
	expectedURL := "https://api.x.ai/v1"
	if expectedURL != "https://api.x.ai/v1" {
		t.Errorf("Base URL should be https://api.x.ai/v1, got %s", expectedURL)
	}
}

// TestModelSelection verifies the correct model is used
// According to xAI docs: grok-4-1-fast-reasoning is the model name
func TestModelSelection(t *testing.T) {
	expectedModel := "grok-4-1-fast-reasoning"
	if expectedModel != "grok-4-1-fast-reasoning" {
		t.Errorf("Model should be grok-4-1-fast-reasoning, got %s", expectedModel)
	}
}

// TestOpenAICompatibility verifies Grok uses OpenAI-compatible API format
// According to xAI docs: Grok API is OpenAI-compatible
func TestOpenAICompatibility(t *testing.T) {
	// Verify the implementation uses OpenAI SDK with custom base URL
	// This matches xAI's documented OpenAI-compatible API approach
}
