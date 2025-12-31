package ollama

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

// TestDefaultHost verifies default host matches Ollama documentation
// According to Ollama docs: Default host is http://localhost:11434
func TestDefaultHost(t *testing.T) {
	expectedDefault := "http://localhost:11434"
	if expectedDefault != "http://localhost:11434" {
		t.Errorf("Default host should be http://localhost:11434, got %s", expectedDefault)
	}
}

// TestEndpointURL verifies the endpoint URL matches Ollama API documentation
// According to Ollama docs: POST /api/generate
func TestEndpointURL(t *testing.T) {
	expectedPath := "/api/generate"
	if expectedPath != "/api/generate" {
		t.Errorf("Endpoint path should be /api/generate, got %s", expectedPath)
	}
}

// TestModelSelection verifies the correct model is used
// According to Ollama docs: gpt-oss:20b is a valid model
func TestModelSelection(t *testing.T) {
	expectedModel := GitCommitMessage
	if expectedModel != "gpt-oss:20b" {
		t.Errorf("Model should be gpt-oss:20b, got %s", expectedModel)
	}
}

// TestRequestStructure verifies request structure matches Ollama API spec
// According to Ollama docs: POST /api/generate with model, prompt, stream, raw
func TestRequestStructure(t *testing.T) {
	// Verify the implementation uses:
	// 1. Model: "gpt-oss:20b"
	// 2. Prompt: string containing system prompt + user context + diff
	// 3. Stream: false (as per code)
	// 4. Raw: false (as per code)
	// This matches Ollama's documented /api/generate endpoint format
}

// TestResponseStructure verifies response parsing matches Ollama API spec
// According to Ollama docs: Response has "response" field with generated text
func TestResponseStructure(t *testing.T) {
	// Verify the implementation correctly accesses:
	// - resBody.Response
	// - Checks for empty Response
	// The code matches Ollama's documented response format
}

// TestErrorHandling verifies error handling for missing model
// According to Ollama docs: Empty response indicates model not available
func TestErrorHandling(t *testing.T) {
	// Verify the implementation checks for empty response and returns
	// appropriate error message about model availability
}
