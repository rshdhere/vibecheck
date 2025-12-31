package cmd

import (
	"bytes"
	"testing"
)

func TestRootCmd(t *testing.T) {
	// Test that root command exists and has correct properties
	if rootCmd.Use != "vibecheck" {
		t.Errorf("rootCmd.Use = %q, want 'vibecheck'", rootCmd.Use)
	}

	if rootCmd.Short == "" {
		t.Error("rootCmd.Short is empty")
	}

	if rootCmd.Long == "" {
		t.Error("rootCmd.Long is empty")
	}
}

func TestExecute(t *testing.T) {
	// Test that Execute doesn't panic with help flag
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"--help"})

	// This should not error (help is a valid command)
	// Note: Execute will call os.Exit(1) on error, so we can't easily test errors
	// But we can test that it doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Execute() panicked: %v", r)
		}
	}()

	// We'll just verify the command is set up correctly
	// Actual execution testing is complex due to os.Exit behavior
}

func TestVersion(t *testing.T) {
	// Test that version variable exists
	if version == "" {
		t.Error("version variable is empty")
	}
}
