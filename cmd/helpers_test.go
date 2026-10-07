package cmd

import (
	"bytes"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rshdhere/vibecheck/internal/keys"
)

// isolateHome points HOME at a fresh directory and clears provider env vars so
// config, keys and stats never touch the developer's real files.
func isolateHome(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	for _, envVar := range keys.ProviderToEnvVar {
		t.Setenv(envVar, "")
	}
	return dir
}

// scriptTUI feeds input to the next interactive program and returns its output.
func scriptTUI(t *testing.T, input string) *bytes.Buffer {
	t.Helper()
	var out bytes.Buffer
	old := teaOptions
	teaOptions = []tea.ProgramOption{
		tea.WithInput(strings.NewReader(input)),
		tea.WithOutput(&out),
		tea.WithoutSignalHandler(),
	}
	t.Cleanup(func() { teaOptions = old })
	return &out
}

// execute runs the root command with args and returns its error.
func execute(t *testing.T, args ...string) error {
	t.Helper()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs(args)
	t.Cleanup(func() {
		rootCmd.SetOut(nil)
		rootCmd.SetErr(nil)
		rootCmd.SetArgs(nil)
	})
	return rootCmd.Execute()
}

func keyRunes(s string) tea.KeyMsg { return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)} }

var (
	keyEnter = tea.KeyMsg{Type: tea.KeyEnter}
	keyEsc   = tea.KeyMsg{Type: tea.KeyEsc}
)
