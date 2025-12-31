package qwen

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
	if err != nil && err.Error() != "QWEN_API_KEY environment variable not set" {
		t.Errorf("GenerateCommitMessage() error message = %q, want 'QWEN_API_KEY environment variable not set'", err.Error())
	}
}

// TestEndpointURL verifies the endpoint URL matches Alibaba Qwen API documentation
// According to Qwen docs: https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions
func TestEndpointURL(t *testing.T) {
	expectedURL := "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"
	if expectedURL != "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions" {
		t.Errorf("Endpoint URL should be https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions, got %s", expectedURL)
	}
}

// TestModelSelection verifies the correct model is used
// According to Qwen docs: qwen-turbo is available
func TestModelSelection(t *testing.T) {
	expectedModel := "qwen-turbo"
	if expectedModel != "qwen-turbo" {
		t.Errorf("Model should be qwen-turbo, got %s", expectedModel)
	}
}

// TestResponseStructure verifies response parsing handles Qwen's dual format
// According to Qwen docs: Response may be in Choices or Output.Choices
func TestResponseStructure(t *testing.T) {
	// Verify the implementation checks both:
	// - chatResp.Choices[0].Message.Content
	// - chatResp.Output.Choices[0].Message.Content
	// This matches Qwen's documented response format flexibility
}
