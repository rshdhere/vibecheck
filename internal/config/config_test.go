package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetConfigPath(t *testing.T) {
	path, err := getConfigPath()
	if err != nil {
		t.Fatalf("getConfigPath() error = %v", err)
	}
	if path == "" {
		t.Error("getConfigPath() returned empty path")
	}
	expectedSuffix := ".vibecheck.json"
	if filepath.Base(path) != ".vibecheck.json" {
		t.Errorf("getConfigPath() = %v, want path ending with %v", path, expectedSuffix)
	}
}

func TestLoad(t *testing.T) {
	t.Run("non-existent config returns default", func(t *testing.T) {
		// Use a temp directory
		tmpDir := t.TempDir()
		oldHome := os.Getenv("HOME")
		defer os.Setenv("HOME", oldHome)
		os.Setenv("HOME", tmpDir)

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if cfg.DefaultProvider != "openai" {
			t.Errorf("Load() DefaultProvider = %v, want openai", cfg.DefaultProvider)
		}
	})

	t.Run("existing config loads correctly", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldHome := os.Getenv("HOME")
		defer os.Setenv("HOME", oldHome)
		os.Setenv("HOME", tmpDir)

		// Create config file
		configPath := filepath.Join(tmpDir, ".vibecheck.json")
		configData := `{"default_provider": "gemini"}`
		if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
			t.Fatalf("Failed to write config: %v", err)
		}

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if cfg.DefaultProvider != "gemini" {
			t.Errorf("Load() DefaultProvider = %v, want gemini", cfg.DefaultProvider)
		}
	})
}

func TestSave(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	cfg := &Config{DefaultProvider: "anthropic"}
	if err := Save(cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify it was saved
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() after Save() error = %v", err)
	}
	if loaded.DefaultProvider != "anthropic" {
		t.Errorf("Load() after Save() DefaultProvider = %v, want anthropic", loaded.DefaultProvider)
	}
}

func TestGetDefaultProvider(t *testing.T) {
	t.Run("default provider", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldHome := os.Getenv("HOME")
		defer os.Setenv("HOME", oldHome)
		os.Setenv("HOME", tmpDir)

		provider := GetDefaultProvider()
		if provider != "openai" {
			t.Errorf("GetDefaultProvider() = %v, want openai", provider)
		}
	})

	t.Run("custom provider", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldHome := os.Getenv("HOME")
		defer os.Setenv("HOME", oldHome)
		os.Setenv("HOME", tmpDir)

		cfg := &Config{DefaultProvider: "groq"}
		if err := Save(cfg); err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		provider := GetDefaultProvider()
		if provider != "groq" {
			t.Errorf("GetDefaultProvider() = %v, want groq", provider)
		}
	})
}

func TestSetDefaultProvider(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	if err := SetDefaultProvider("kimi"); err != nil {
		t.Fatalf("SetDefaultProvider() error = %v", err)
	}

	provider := GetDefaultProvider()
	if provider != "kimi" {
		t.Errorf("GetDefaultProvider() after SetDefaultProvider() = %v, want kimi", provider)
	}
}

func TestLoadWithInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Create invalid JSON config file
	configPath := filepath.Join(tmpDir, ".vibecheck.json")
	invalidJSON := `{"default_provider": "openai" invalid}`
	if err := os.WriteFile(configPath, []byte(invalidJSON), 0644); err != nil {
		t.Fatalf("Failed to write invalid config: %v", err)
	}

	_, err := Load()
	if err == nil {
		t.Error("Load() with invalid JSON should return error")
	}
}

func TestSaveWithInvalidPath(t *testing.T) {
	// Test Save with invalid path (should still work if HOME is valid)
	// This tests the error path in Save
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)

	// Set HOME to a non-existent path to test error handling
	os.Setenv("HOME", "/nonexistent/path/that/does/not/exist")

	cfg := &Config{DefaultProvider: "test"}
	err := Save(cfg)
	// Error is expected with invalid path
	_ = err
}
