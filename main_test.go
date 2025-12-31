package main

import (
	"os"
	"testing"
)

func TestIsVersionOrHelp(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected bool
	}{
		{
			name:     "version flag",
			args:     []string{"--version"},
			expected: true,
		},
		{
			name:     "short version flag",
			args:     []string{"-v"},
			expected: true,
		},
		{
			name:     "version command",
			args:     []string{"version"},
			expected: true,
		},
		{
			name:     "help flag",
			args:     []string{"--help"},
			expected: true,
		},
		{
			name:     "short help flag",
			args:     []string{"-h"},
			expected: true,
		},
		{
			name:     "help command",
			args:     []string{"help"},
			expected: true,
		},
		{
			name:     "no args",
			args:     []string{},
			expected: true,
		},
		{
			name:     "other command",
			args:     []string{"commit"},
			expected: false,
		},
		{
			name:     "version in middle",
			args:     []string{"commit", "--version"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()
			os.Args = append([]string{"vibecheck"}, tt.args...)
			result := isVersionOrHelp()
			if result != tt.expected {
				t.Errorf("isVersionOrHelp() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestLoadDotEnvIfPresent(t *testing.T) {
	// Test with non-existent .env file (should not error)
	t.Run("no .env file", func(t *testing.T) {
		// Create a temp directory without .env
		tmpDir := t.TempDir()
		oldDir, _ := os.Getwd()
		defer os.Chdir(oldDir)
		os.Chdir(tmpDir)
		loadDotEnvIfPresent() // Should not panic or error
	})

	// Test with existing .env file
	t.Run("existing .env file", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldDir, _ := os.Getwd()
		defer os.Chdir(oldDir)
		os.Chdir(tmpDir)

		// Create a valid .env file
		envContent := "TEST_KEY=test_value\n"
		err := os.WriteFile(".env", []byte(envContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create .env file: %v", err)
		}

		loadDotEnvIfPresent() // Should not panic
	})
}
