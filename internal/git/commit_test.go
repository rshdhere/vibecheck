package git_test

import (
	"testing"
)

func TestCommitWMessage(t *testing.T) {
	// Skip this test as git commit -em opens an editor which will hang in tests
	// The function requires interactive input which is not suitable for automated testing
	// In a real scenario, this would require mocking the exec.Command or using a different approach
	t.Skip("Skipping TestCommitWMessage - git commit -em requires interactive editor input")
}
