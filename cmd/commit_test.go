package cmd

import (
	"testing"
)

func TestDetectMissingEnvVar(t *testing.T) {
	tests := []struct {
		name     string
		errMsg   string
		expected string
	}{
		{
			name:     "valid env var error",
			errMsg:   "OPENAI_API_KEY environment variable not set",
			expected: "OPENAI_API_KEY",
		},
		{
			name:     "different env var",
			errMsg:   "GEMINI_API_KEY environment variable not set",
			expected: "GEMINI_API_KEY",
		},
		{
			name:     "no suffix match",
			errMsg:   "some other error",
			expected: "",
		},
		{
			name:     "empty error",
			errMsg:   "",
			expected: "",
		},
		{
			name:     "partial match",
			errMsg:   "environment variable not set",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &testError{msg: tt.errMsg}
			result := detectMissingEnvVar(err)
			if result != tt.expected {
				t.Errorf("detectMissingEnvVar() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestDetectMissingModel(t *testing.T) {
	tests := []struct {
		name     string
		errMsg   string
		expected string
	}{
		{
			name:     "valid model error",
			errMsg:   "model 'gpt-4o-mini' not found",
			expected: "gpt-4o-mini",
		},
		{
			name:     "different model",
			errMsg:   "model 'claude-3.5-haiku' not found",
			expected: "claude-3.5-haiku",
		},
		{
			name:     "no model pattern",
			errMsg:   "some other error",
			expected: "",
		},
		{
			name:     "empty error",
			errMsg:   "",
			expected: "",
		},
		{
			name:     "partial match",
			errMsg:   "model 'test",
			expected: "",
		},
		{
			name:     "missing closing quote",
			errMsg:   "model 'test not found",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &testError{msg: tt.errMsg}
			result := detectMissingModel(err)
			if result != tt.expected {
				t.Errorf("detectMissingModel() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// testError is a simple error implementation for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
