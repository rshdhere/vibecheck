package cmd

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rshdhere/vibecheck/internal/llm"
	"github.com/rshdhere/vibecheck/internal/stats"
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

// fakeProvider returns a canned message or error.
type fakeProvider struct {
	msg string
	err error
}

func (f fakeProvider) GenerateCommitMessage(context.Context, string, string) (string, error) {
	return f.msg, f.err
}

func init() {
	llm.Register("test-ok", fakeProvider{msg: "feat(cmd): add thing\n\n- detail"})
	llm.Register("test-missing-key", fakeProvider{err: errors.New("FAKE_API_KEY environment variable not set")})
	llm.Register("test-missing-model", fakeProvider{err: errors.New("model 'fake-model' not found")})
	llm.Register("test-fail", fakeProvider{err: errors.New("boom")})
}

// setupRepo creates a git repo with one staged file and makes it the working directory.
func setupRepo(t *testing.T, stage bool) string {
	t.Helper()
	isolateHome(t)
	dir := t.TempDir()
	t.Chdir(dir)
	t.Setenv("GIT_EDITOR", "true")
	runGit(t, "init", "-q")
	runGit(t, "config", "user.email", "test@example.com")
	runGit(t, "config", "user.name", "Test User")
	runGit(t, "config", "commit.gpgsign", "false")
	if stage {
		if err := os.WriteFile(filepath.Join(dir, "file.txt"), []byte("hello\n"), 0644); err != nil {
			t.Fatal(err)
		}
		runGit(t, "add", "file.txt")
	}
	return dir
}

func runGit(t *testing.T, args ...string) string {
	t.Helper()
	out, err := exec.Command("git", args...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(out)
}

func headMessage(t *testing.T) string {
	t.Helper()
	out, err := exec.Command("git", "log", "-1", "--format=%B").CombinedOutput()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func TestCommitCmdCommitsGeneratedMessage(t *testing.T) {
	setupRepo(t, true)

	if err := execute(t, "commit", "--provider", "test-ok", "--prompt", "ctx"); err != nil {
		t.Fatalf("commit error = %v", err)
	}

	if got := headMessage(t); got != "feat(cmd): add thing\n\n- detail" {
		t.Errorf("HEAD message = %q, want the generated message", got)
	}
	recent, err := stats.GetRecentCommits(1)
	if err != nil || len(recent) != 1 {
		t.Fatalf("recent commits = %v, %v; want one record", recent, err)
	}
	if recent[0].Model != "test-ok" || recent[0].CommitMsg != "feat(cmd): add thing" {
		t.Errorf("recorded %+v, want model test-ok and the first message line", recent[0])
	}
}

func TestCommitCmdHandledProviderErrors(t *testing.T) {
	for _, provider := range []string{"test-missing-key", "test-missing-model"} {
		t.Run(provider, func(t *testing.T) {
			setupRepo(t, true)

			if err := execute(t, "commit", "--provider", provider, "--prompt", ""); err != nil {
				t.Fatalf("commit error = %v, want nil after showing a notice", err)
			}
			if got := headMessage(t); got != "" {
				t.Errorf("HEAD message = %q, want no commit", got)
			}
		})
	}
}

func TestCommitCmdNothingStaged(t *testing.T) {
	setupRepo(t, false)

	if err := execute(t, "commit", "--provider", "test-ok", "--prompt", ""); err != nil {
		t.Fatalf("commit error = %v, want nil after the stage reminder", err)
	}
	if got := headMessage(t); got != "" {
		t.Errorf("HEAD message = %q, want no commit", got)
	}
}

func TestCommitCmdErrors(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		setup    func(t *testing.T)
		wantErr  string
	}{
		{"not a git repo", "test-ok", func(t *testing.T) { isolateHome(t); t.Chdir(t.TempDir()) }, "staged changes"},
		{"unknown provider", "no-such-provider", func(t *testing.T) { setupRepo(t, true) }, llm.ErrNoProvider.Error()},
		{"provider failure", "test-fail", func(t *testing.T) { setupRepo(t, true) }, "generated commit message: boom"},
		{"git commit aborted", "test-ok", func(t *testing.T) { setupRepo(t, true); t.Setenv("GIT_EDITOR", "false") }, "commit with message"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)

			err := execute(t, "commit", "--provider", tt.provider, "--prompt", "")
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("commit error = %v, want it to contain %q", err, tt.wantErr)
			}
		})
	}
}
