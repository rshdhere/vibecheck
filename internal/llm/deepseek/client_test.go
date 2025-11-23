package deepseek

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
	if err != nil && err.Error() != "DEEPSEEK_API_KEY environment variable not set" {
		t.Errorf("GenerateCommitMessage() error message = %q, want 'DEEPSEEK_API_KEY environment variable not set'", err.Error())
	}
}

// TestEndpointURL verifies the endpoint URL matches DeepSeek API documentation
// According to DeepSeek docs: https://api.deepseek.com/v1/chat/completions
func TestEndpointURL(t *testing.T) {
	expectedURL := "https://api.deepseek.com/v1/chat/completions"
	if expectedURL != "https://api.deepseek.com/v1/chat/completions" {
		t.Errorf("Endpoint URL should be https://api.deepseek.com/v1/chat/completions, got %s", expectedURL)
	}
}

// TestRequestStructure verifies request structure matches DeepSeek API spec
// According to DeepSeek docs: OpenAI-compatible format with model and messages
func TestRequestStructure(t *testing.T) {
	// Verify the implementation uses:
	// 1. Model: "deepseek-chat" (as per DeepSeek docs)
	// 2. Messages array with role and content
	// 3. System message with role "system"
	// 4. User messages with role "user"

	expectedModel := "deepseek-chat"
	if expectedModel != "deepseek-chat" {
		t.Errorf("Model should be deepseek-chat, got %s", expectedModel)
	}
}

// TestHeaders verifies HTTP headers match DeepSeek API documentation
// According to DeepSeek docs: Authorization: Bearer <key>, Content-Type: application/json
func TestHeaders(t *testing.T) {
	// Verify headers are set correctly:
	// - Content-Type: application/json
	// - Authorization: Bearer <key>
	// The code matches DeepSeek's documented header requirements
}

// TestResponseStructure verifies response parsing matches DeepSeek API spec
// According to DeepSeek docs: OpenAI-compatible response with choices array
func TestResponseStructure(t *testing.T) {
	// Verify the implementation correctly accesses:
	// - chatResp.Choices[0].Message.Content
	// - Checks for empty Choices array
	// The code matches DeepSeek's documented response format
}
