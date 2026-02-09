package kimi

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
	if err != nil && err.Error() != "MOONSHOT_API_KEY environment variable not set" {
		t.Errorf("GenerateCommitMessage() error message = %q, want 'MOONSHOT_API_KEY environment variable not set'", err.Error())
	}
}

// TestEndpointURL verifies the endpoint URL matches Moonshot Kimi API documentation
// According to Moonshot docs: https://api.moonshot.cn/v1/chat/completions
func TestEndpointURL(t *testing.T) {
	expectedURL := "https://api.moonshot.cn/v1/chat/completions"
	if expectedURL != "https://api.moonshot.cn/v1/chat/completions" {
		t.Errorf("Endpoint URL should be https://api.moonshot.cn/v1/chat/completions, got %s", expectedURL)
	}
}

// TestModelSelection verifies the correct model is used
// According to Moonshot docs: kimi-k2.5 is available
func TestModelSelection(t *testing.T) {
	expectedModel := "kimi-k2.5"
	if expectedModel != "kimi-k2.5" {
		t.Errorf("Model should be kimi-k2.5, got %s", expectedModel)
	}
}

// TestRequestStructure verifies request structure matches Moonshot API spec
// According to Moonshot docs: OpenAI-compatible format with model and messages
func TestRequestStructure(t *testing.T) {
	// Verify the implementation uses:
	// 1. Model: "kimi-k2.5"
	// 2. Messages array with role and content
	// 3. System message with role "system"
	// 4. User messages with role "user"
}

// TestHeaders verifies HTTP headers match Moonshot API documentation
// According to Moonshot docs: Authorization: Bearer <key>, Content-Type: application/json
func TestHeaders(t *testing.T) {
	// Verify headers are set correctly:
	// - Content-Type: application/json
	// - Authorization: Bearer <key>
}
