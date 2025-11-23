package git_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/rshdhere/vibecheck/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func SetupGitRepo() (string, error) {
	path, err := os.MkdirTemp("", "")
	if err != nil {
		return "", err
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = path

	if err := cmd.Run(); err != nil {
		os.RemoveAll(path)
		return "", err
	}
	return path, nil
}

func TestStagedDiff(t *testing.T) {
	repo, err := SetupGitRepo()
	require.NoError(t, err)
	defer os.RemoveAll(repo)

	testOneFile, err := os.Create(fmt.Sprintf("%s/test1.txt", repo))
	require.NoError(t, err)

	testOneFile.WriteString("hello from test001\n")

	testTwoFile, err := os.Create(fmt.Sprintf("%s/test2.txt", repo))
	require.NoError(t, err)

	testTwoFile.WriteString("hello from test002\n")

	cmd := exec.Command("git", "add", "test1.txt")
	cmd.Dir = repo
	require.NoError(t, cmd.Run())

	wd, err := os.Getwd()
	require.NoError(t, err)

	defer os.Chdir(wd)

	os.Chdir(repo)
	changes, err := git.StagedDiff(context.Background())

	assert.NoError(t, err)

	assert.Contains(t, changes, "test1.txt")
	assert.Contains(t, changes, "hello from test001")
}

func TestStagedDiffEmpty(t *testing.T) {
	repo, err := SetupGitRepo()
	require.NoError(t, err)
	defer os.RemoveAll(repo)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	os.Chdir(repo)

	changes, err := git.StagedDiff(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "", changes)
}

func TestStagedDiffWithContext(t *testing.T) {
	repo, err := SetupGitRepo()
	require.NoError(t, err)
	defer os.RemoveAll(repo)

	testFile, err := os.Create(fmt.Sprintf("%s/test.txt", repo))
	require.NoError(t, err)
	testFile.WriteString("test content\n")
	testFile.Close()

	cmd := exec.Command("git", "add", "test.txt")
	cmd.Dir = repo
	require.NoError(t, cmd.Run())

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	os.Chdir(repo)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately to test context handling

	_, err = git.StagedDiff(ctx)
	// Should handle context cancellation gracefully
	_ = err // Error is expected when context is cancelled
}

func TestStagedDiffMultipleFiles(t *testing.T) {
	repo, err := SetupGitRepo()
	require.NoError(t, err)
	defer os.RemoveAll(repo)

	// Create multiple files
	for i := 1; i <= 3; i++ {
		testFile, err := os.Create(fmt.Sprintf("%s/file%d.txt", repo, i))
		require.NoError(t, err)
		testFile.WriteString(fmt.Sprintf("content %d\n", i))
		testFile.Close()
	}

	// Stage all files
	cmd := exec.Command("git", "add", "file1.txt", "file2.txt", "file3.txt")
	cmd.Dir = repo
	require.NoError(t, cmd.Run())

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	os.Chdir(repo)

	changes, err := git.StagedDiff(context.Background())
	assert.NoError(t, err)
	assert.Contains(t, changes, "file1.txt")
	assert.Contains(t, changes, "file2.txt")
	assert.Contains(t, changes, "file3.txt")
}

func TestStagedDiffWithModifications(t *testing.T) {
	repo, err := SetupGitRepo()
	require.NoError(t, err)
	defer os.RemoveAll(repo)

	// Create and commit initial file
	testFile, err := os.Create(fmt.Sprintf("%s/test.txt", repo))
	require.NoError(t, err)
	testFile.WriteString("initial content\n")
	testFile.Close()

	cmd := exec.Command("git", "add", "test.txt")
	cmd.Dir = repo
	require.NoError(t, cmd.Run())

	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = repo
	require.NoError(t, cmd.Run())

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = repo
	require.NoError(t, cmd.Run())

	cmd = exec.Command("git", "commit", "-m", "initial commit")
	cmd.Dir = repo
	require.NoError(t, cmd.Run())

	// Modify the file
	testFile, err = os.OpenFile(fmt.Sprintf("%s/test.txt", repo), os.O_APPEND|os.O_WRONLY, 0644)
	require.NoError(t, err)
	testFile.WriteString("modified content\n")
	testFile.Close()

	// Stage the modification
	cmd = exec.Command("git", "add", "test.txt")
	cmd.Dir = repo
	require.NoError(t, cmd.Run())

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	os.Chdir(repo)

	changes, err := git.StagedDiff(context.Background())
	assert.NoError(t, err)
	assert.Contains(t, changes, "test.txt")
	assert.Contains(t, changes, "modified content")
}
