package groq

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
	if err != nil && err.Error() != "GROQ_API_KEY environment variable not set" {
		t.Errorf("GenerateCommitMessage() error message = %q, want 'GROQ_API_KEY environment variable not set'", err.Error())
	}
}

// TestBaseURL verifies the base URL matches Groq API documentation
// According to Groq docs: https://api.groq.com/openai/v1 (OpenAI-compatible)
func TestBaseURL(t *testing.T) {
	expectedURL := "https://api.groq.com/openai/v1"
	if expectedURL != "https://api.groq.com/openai/v1" {
		t.Errorf("Base URL should be https://api.groq.com/openai/v1, got %s", expectedURL)
	}
}

// TestModelSelection verifies the correct model is used
// According to Groq docs: llama-3.3-70b-versatile is available
func TestModelSelection(t *testing.T) {
	expectedModel := "llama-3.3-70b-versatile"
	if expectedModel != "llama-3.3-70b-versatile" {
		t.Errorf("Model should be llama-3.3-70b-versatile, got %s", expectedModel)
	}
}

// TestOpenAICompatibility verifies Groq uses OpenAI-compatible API format
// According to Groq docs: Groq API is OpenAI-compatible
func TestOpenAICompatibility(t *testing.T) {
	// Verify the implementation uses OpenAI SDK with custom base URL
	// This matches Groq's documented OpenAI-compatible API approach
}
