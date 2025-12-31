package git

import "testing"

func FuzzCommitMessage(f *testing.F) {
	f.Add("feat: initial commit")
	f.Add("fix(core): handle nil pointer")
	f.Add("🚀 release v1.0.0")
	f.Add("")

	f.Fuzz(func(t *testing.T, msg string) {
		// CALL REAL FUNCTION HERE
		_ = ParseCommit(msg) // example — replace with real one
	})
}
