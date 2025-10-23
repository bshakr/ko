// Package git provides utilities for working with git repositories and worktrees.
//
// This package wraps git commands to:
//   - Check if the current directory is in a git repository
//   - Create and remove git worktrees
//   - Detect if running inside a worktree vs. main repository
//   - Get repository and worktree paths
//
// All operations that may be long-running support context-based cancellation,
// allowing users to interrupt operations with Ctrl+C.
//
// Worktree Detection:
// The package can distinguish between the main repository and worktrees by
// comparing git-dir and git-common-dir. This is essential for ko's functionality
// since configuration is stored in the main repository, not individual worktrees.
package git

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// IsGitRepo checks if the current directory is in a git repository
func IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	return err == nil
}

// CreateWorktree creates a new git worktree at the specified path
func CreateWorktree(path string) error {
	cmd := exec.Command("git", "worktree", "add", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}
	return nil
}

// CreateWorktreeWithContext creates a new git worktree at the specified path with cancellation support
func CreateWorktreeWithContext(ctx context.Context, path string) error {
	cmd := exec.CommandContext(ctx, "git", "worktree", "add", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.Canceled {
			return fmt.Errorf("operation cancelled")
		}
		return fmt.Errorf("%s", string(output))
	}
	return nil
}

// RemoveWorktree removes a git worktree at the specified path
func RemoveWorktree(path string) error {
	cmd := exec.Command("git", "worktree", "remove", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}
	return nil
}

// RemoveWorktreeWithContext removes a git worktree at the specified path with cancellation support
func RemoveWorktreeWithContext(ctx context.Context, path string) error {
	cmd := exec.CommandContext(ctx, "git", "worktree", "remove", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.Canceled {
			return fmt.Errorf("operation cancelled")
		}
		return fmt.Errorf("%s", string(output))
	}
	return nil
}

// GetRepoName returns the name of the current git repository
func GetRepoName() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get repository root: %w", err)
	}

	topLevel := strings.TrimSpace(string(output))
	repoName := filepath.Base(topLevel)
	return repoName, nil
}

// IsInWorktree checks if the current directory is inside a worktree
func IsInWorktree() bool {
	gitDirCmd := exec.Command("git", "rev-parse", "--git-dir")
	gitDirOutput, err := gitDirCmd.Output()
	if err != nil {
		return false
	}

	commonDirCmd := exec.Command("git", "rev-parse", "--git-common-dir")
	commonDirOutput, err := commonDirCmd.Output()
	if err != nil {
		return false
	}

	gitDir := strings.TrimSpace(string(gitDirOutput))
	commonDir := strings.TrimSpace(string(commonDirOutput))

	// If git-dir and git-common-dir are different, we're in a worktree
	return gitDir != commonDir
}

// GetMainRepoRoot returns the root of the main repository (not the worktree)
func GetMainRepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--git-common-dir")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git common dir: %w", err)
	}

	commonDir := strings.TrimSpace(string(output))
	// The common dir is .git, so we need to go up one level
	mainRepoRoot := filepath.Dir(commonDir)
	return mainRepoRoot, nil
}

// GetCurrentWorktreePath returns the current worktree path
func GetCurrentWorktreePath() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get worktree path: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}
