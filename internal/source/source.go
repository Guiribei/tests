package source

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func Sync(ctx context.Context, directory, repositoryURL string) error {
	if _, err := os.Stat(filepath.Join(directory, ".git")); err == nil {
		cmd := exec.CommandContext(ctx, "git", "-C", directory, "pull", "--ff-only")
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("git pull: %w: %s", err, output)
		}
		return nil
	}

	entries, err := os.ReadDir(directory)
	if err == nil && len(entries) > 0 {
		return fmt.Errorf("%s não é um checkout Git e não está vazio", directory)
	}
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(directory), 0o755); err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, "git", "clone", "--depth", "1", repositoryURL, directory)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git clone: %w: %s", err, output)
	}
	return nil
}
