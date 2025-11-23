package keys

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetKeysPath(t *testing.T) {
	path, err := getKeysPath()
	if err != nil {
		t.Fatalf("getKeysPath() error = %v", err)
	}
	if path == "" {
		t.Error("getKeysPath() returned empty path")
	}
	if filepath.Base(path) != ".vibecheck_keys.json" {
		t.Errorf("getKeysPath() = %v, want path ending with .vibecheck_keys.json", path)
	}
}

func TestLoad(t *testing.T) {
	t.Run("non-existent keys returns empty", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldHome := os.Getenv("HOME")
		defer os.Setenv("HOME", oldHome)
		os.Setenv("HOME", tmpDir)

		keys, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if keys == nil {
			t.Error("Load() returned nil")
		}
	})

	t.Run("existing keys loads correctly", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldHome := os.Getenv("HOME")
		defer os.Setenv("HOME", oldHome)
		os.Setenv("HOME", tmpDir)

		keysPath := filepath.Join(tmpDir, ".vibecheck_keys.json")
		keysData := `{"openai": "sk-test123", "gemini": "test-gemini-key"}`
		if err := os.WriteFile(keysPath, []byte(keysData), 0600); err != nil {
			t.Fatalf("Failed to write keys: %v", err)
		}

		keys, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if keys.OpenAI != "sk-test123" {
			t.Errorf("Load() OpenAI = %v, want sk-test123", keys.OpenAI)
		}
		if keys.Gemini != "test-gemini-key" {
			t.Errorf("Load() Gemini = %v, want test-gemini-key", keys.Gemini)
		}
	})
}

func TestSave(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	keys := &Keys{
		OpenAI:    "sk-test123",
		Anthropic: "sk-ant-test",
	}
	if err := Save(keys); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify it was saved
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() after Save() error = %v", err)
	}
	if loaded.OpenAI != "sk-test123" {
		t.Errorf("Load() after Save() OpenAI = %v, want sk-test123", loaded.OpenAI)
	}
	if loaded.Anthropic != "sk-ant-test" {
		t.Errorf("Load() after Save() Anthropic = %v, want sk-ant-test", loaded.Anthropic)
	}
}

func TestGetAPIKey(t *testing.T) {
	t.Run("from environment variable", func(t *testing.T) {
		os.Setenv("OPENAI_API_KEY", "env-key-123")
		defer os.Unsetenv("OPENAI_API_KEY")

		key, exists := GetAPIKey("openai")
		if !exists {
			t.Error("GetAPIKey() exists = false, want true")
		}
		if key != "env-key-123" {
			t.Errorf("GetAPIKey() = %v, want env-key-123", key)
		}
	})

	t.Run("from keys file when env not set", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldHome := os.Getenv("HOME")
		defer os.Setenv("HOME", oldHome)
		os.Setenv("HOME", tmpDir)

		keys := &Keys{OpenAI: "file-key-456"}
		if err := Save(keys); err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		key, exists := GetAPIKey("openai")
		if !exists {
			t.Error("GetAPIKey() exists = false, want true")
		}
		if key != "file-key-456" {
			t.Errorf("GetAPIKey() = %v, want file-key-456", key)
		}
	})

	t.Run("not found", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldHome := os.Getenv("HOME")
		defer os.Setenv("HOME", oldHome)
		os.Setenv("HOME", tmpDir)

		key, exists := GetAPIKey("nonexistent")
		if exists {
			t.Error("GetAPIKey() exists = true, want false")
		}
		if key != "" {
			t.Errorf("GetAPIKey() = %v, want empty string", key)
		}
	})
}

func TestSetAPIKey(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	tests := []struct {
		provider string
		key      string
	}{
		{"openai", "sk-test-openai"},
		{"gemini", "test-gemini"},
		{"anthropic", "sk-ant-test"},
		{"groq", "gsk-test"},
		{"grok", "xai-test"},
		{"kimi", "moonshot-test"},
		{"qwen", "qwen-test"},
		{"deepseek", "deepseek-test"},
		{"perplexity", "pplx-test"},
		{"ollama", "http://localhost:11434"},
	}

	for _, tt := range tests {
		t.Run(tt.provider, func(t *testing.T) {
			if err := SetAPIKey(tt.provider, tt.key); err != nil {
				t.Fatalf("SetAPIKey(%v, %v) error = %v", tt.provider, tt.key, err)
			}

			key, exists := GetAPIKey(tt.provider)
			if !exists {
				t.Errorf("GetAPIKey(%v) exists = false after SetAPIKey", tt.provider)
			}
			if key != tt.key {
				t.Errorf("GetAPIKey(%v) = %v, want %v", tt.provider, key, tt.key)
			}
		})
	}

	t.Run("unknown provider", func(t *testing.T) {
		err := SetAPIKey("unknown", "test")
		if err == nil {
			t.Error("SetAPIKey() with unknown provider should return error")
		}
	})
}

func TestGetAllKeys(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	keys := &Keys{
		OpenAI:     "sk-test123456789",
		Gemini:     "gemini-key-123456",
		OllamaHost: "http://localhost:11434",
	}
	if err := Save(keys); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	allKeys, err := GetAllKeys()
	if err != nil {
		t.Fatalf("GetAllKeys() error = %v", err)
	}

	if allKeys["openai"] != "sk-t...6789" {
		t.Errorf("GetAllKeys() openai = %v, want masked key", allKeys["openai"])
	}
	if allKeys["gemini"] != "gemi...3456" {
		t.Errorf("GetAllKeys() gemini = %v, want masked key", allKeys["gemini"])
	}
	if allKeys["ollama"] != "http://localhost:11434" {
		t.Errorf("GetAllKeys() ollama = %v, want http://localhost:11434 (not masked)", allKeys["ollama"])
	}
}

func TestMaskKey(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected string
	}{
		{"long key", "sk-test123456789", "sk-t...6789"},
		{"short key", "short", "****"},
		{"very short key", "abc", "****"},
		{"exactly 8 chars", "12345678", "****"},
		{"9 chars", "123456789", "1234...6789"},
		{"empty key", "", "****"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskKey(tt.key)
			if result != tt.expected {
				t.Errorf("maskKey(%v) = %v, want %v", tt.key, result, tt.expected)
			}
		})
	}
}

func TestGetAllKeysWithAllProviders(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	keys := &Keys{
		OpenAI:     "sk-openai-123456789",
		Gemini:     "gemini-key-123456",
		Anthropic:  "sk-ant-123456789",
		Groq:       "gsk-groq-123456789",
		Grok:       "xai-grok-123456789",
		Kimi:       "moonshot-kimi-123456789",
		Qwen:       "qwen-key-123456789",
		DeepSeek:   "deepseek-key-123456789",
		Perplexity: "pplx-key-123456789",
		OllamaHost: "http://localhost:11434",
	}
	if err := Save(keys); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	allKeys, err := GetAllKeys()
	if err != nil {
		t.Fatalf("GetAllKeys() error = %v", err)
	}

	// Verify all providers are present
	expectedProviders := []string{"openai", "gemini", "anthropic", "groq", "grok", "kimi", "qwen", "deepseek", "perplexity", "ollama"}
	for _, provider := range expectedProviders {
		if _, exists := allKeys[provider]; !exists {
			t.Errorf("GetAllKeys() missing provider: %s", provider)
		}
	}
}

func TestLoadWithInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Create invalid JSON keys file
	keysPath := filepath.Join(tmpDir, ".vibecheck_keys.json")
	invalidJSON := `{"openai": "sk-test" invalid}`
	if err := os.WriteFile(keysPath, []byte(invalidJSON), 0600); err != nil {
		t.Fatalf("Failed to write invalid keys: %v", err)
	}

	_, err := Load()
	if err == nil {
		t.Error("Load() with invalid JSON should return error")
	}
}
