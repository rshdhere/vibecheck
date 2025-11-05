// Package git is used in here to autmomate commit messages coming from Ollama
package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

func CommitWMessage(ctx context.Context, msg string) error {
	// testing what the AI does for this one
	// and this one too
	cmd := exec.CommandContext(ctx, "git", "commit", "-em", msg)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		return err
	}
	cmd.Run()
	return nil
}
