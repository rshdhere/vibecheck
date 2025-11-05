package git

import (
	"context"
	"os/exec"
)

func StagedDiff(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "diff", "--staged")

	res, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(res), nil
}
