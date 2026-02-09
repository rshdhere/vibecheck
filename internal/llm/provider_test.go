package llm

import (
	"context"
	"errors"
	"testing"
)

// mockProvider is a test implementation of Provider
type mockProvider struct {
	generateFunc func(ctx context.Context, diff string, additionalContext string) (string, error)
}

func (m *mockProvider) GenerateCommitMessage(ctx context.Context, diff string, additionalContext string) (string, error) {
	if m.generateFunc != nil {
		return m.generateFunc(ctx, diff, additionalContext)
	}
	return "test commit message", nil
}

func TestRegister(t *testing.T) {
	// Reset providers map
	providers = make(map[string]Provider)

	provider := &mockProvider{}
	Register("test-provider", provider)

	registered := GetRegisteredNames()
	found := false
	for _, name := range registered {
		if name == "test-provider" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Register() did not register provider")
	}
}

func TestGetRegisteredNames(t *testing.T) {
	// Reset providers map
	providers = make(map[string]Provider)

	Register("provider1", &mockProvider{})
	Register("provider2", &mockProvider{})

	names := GetRegisteredNames()
	if len(names) < 2 {
		t.Errorf("GetRegisteredNames() length = %v, want at least 2", len(names))
	}

	found1, found2 := false, false
	for _, name := range names {
		if name == "provider1" {
			found1 = true
		}
		if name == "provider2" {
			found2 = true
		}
	}
	if !found1 || !found2 {
		t.Error("GetRegisteredNames() missing registered providers")
	}
}

func TestGetProvider(t *testing.T) {
	// Reset providers map
	providers = make(map[string]Provider)

	provider := &mockProvider{}
	Register("test-provider", provider)

	got, err := GetProvider("test-provider")
	if err != nil {
		t.Fatalf("GetProvider() error = %v", err)
	}
	if got != provider {
		t.Error("GetProvider() returned wrong provider")
	}

	// Test non-existent provider
	_, err = GetProvider("non-existent")
	if err == nil {
		t.Error("GetProvider() with non-existent provider should return error")
	}
	if !errors.Is(err, ErrNoProvider) {
		t.Errorf("GetProvider() error = %v, want ErrNoProvider", err)
	}
}

func TestMockProvider(t *testing.T) {
	ctx := context.Background()
	provider := &mockProvider{
		generateFunc: func(ctx context.Context, diff string, additionalContext string) (string, error) {
			return "custom message", nil
		},
	}

	msg, err := provider.GenerateCommitMessage(ctx, "test diff", "")
	if err != nil {
		t.Fatalf("GenerateCommitMessage() error = %v", err)
	}
	if msg != "custom message" {
		t.Errorf("GenerateCommitMessage() = %v, want custom message", msg)
	}
}

func TestGetSystemPrompt(t *testing.T) {
	tests := []struct {
		name       string
		provider   string
		wantPrompt string
	}{
		{
			name:       "shared prompt provider",
			provider:   "openai",
			wantPrompt: sharedCommitMessageSystemPrompt,
		},
		{
			name:       "gemini prompt provider",
			provider:   "gemini",
			wantPrompt: geminiCommitMessageSystemPrompt,
		},
		{
			name:       "ollama prompt provider",
			provider:   "ollama",
			wantPrompt: ollamaCommitMessageSystemPrompt,
		},
		{
			name:       "unknown provider falls back to shared prompt",
			provider:   "unknown-provider",
			wantPrompt: sharedCommitMessageSystemPrompt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetSystemPrompt(tt.provider)
			if got != tt.wantPrompt {
				t.Errorf("GetSystemPrompt(%q) returned unexpected prompt", tt.provider)
			}
		})
	}
}
