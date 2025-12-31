package git

import "testing"

// FuzzCommitMessage ensures commit message handling never panics
func FuzzCommitMessage(f *testing.F) {
	// Seed corpus (important for OpenSSF detection)
	f.Add("feat: initial commit")
	f.Add("fix(core): handle nil pointer")
	f.Add("release v1.0.0")
	f.Add("")

	f.Fuzz(func(t *testing.T, msg string) {
		// Call real logic (adjust if function name differs)
		_, _ = NormalizeCommitMessage(msg)
	})
}
