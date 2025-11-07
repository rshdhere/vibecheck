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

	assert.Equal(t, `diff --git a/test1.txt b/test1.txt
new file mode 100644
index 0000000..f119867
--- /dev/null
+++ b/test1.txt
@@ -0,0 +1 @@
+hello from test001
`, changes)
}
