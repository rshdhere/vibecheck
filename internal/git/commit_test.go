package git_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rshdhere/vibecheck/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stagedRepo creates a repo with one staged file, makes it the working directory
// and sets GIT_EDITOR so `git commit -e` runs without an interactive editor.
func stagedRepo(t *testing.T, editor string) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	t.Setenv("GIT_EDITOR", editor)
	dir := t.TempDir()
	t.Chdir(dir)
	for _, args := range [][]string{
		{"init", "-q"},
		{"config", "user.email", "test@example.com"},
		{"config", "user.name", "Test User"},
	} {
		require.NoError(t, exec.Command("git", args...).Run())
	}
	require.NoError(t, os.WriteFile(filepath.Join(dir, "file.txt"), []byte("hello\n"), 0644))
	require.NoError(t, exec.Command("git", "add", "file.txt").Run())
}

func TestCommitWMessage(t *testing.T) {
	stagedRepo(t, "true")

	require.NoError(t, git.CommitWMessage(context.Background(), "feat: add file"))

	out, err := exec.Command("git", "log", "-1", "--format=%s").Output()
	require.NoError(t, err)
	assert.Equal(t, "feat: add file", strings.TrimSpace(string(out)))
}

func TestCommitWMessageEditorAborts(t *testing.T) {
	stagedRepo(t, "false")

	assert.Error(t, git.CommitWMessage(context.Background(), "feat: add file"))
}
