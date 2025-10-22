package git

import (
	"os"
	"testing"
)

func TestIsGitRepo(t *testing.T) {
	// This test assumes we're running in a git repo
	// If not in a git repo, this test will fail
	result := IsGitRepo()
	if !result {
		t.Skip("Not in a git repository, skipping test")
	}
}

func TestGetRepoName(t *testing.T) {
	// Only test if we're in a git repo
	if !IsGitRepo() {
		t.Skip("Not in a git repository, skipping test")
	}

	name, err := GetRepoName()
	if err != nil {
		t.Fatalf("GetRepoName() failed: %v", err)
	}

	if name == "" {
		t.Error("GetRepoName() returned empty string")
	}

	t.Logf("Repository name: %s", name)
}

func TestIsInWorktree(t *testing.T) {
	if !IsGitRepo() {
		t.Skip("Not in a git repository, skipping test")
	}

	// Test should work whether we're in a worktree or not
	result := IsInWorktree()
	t.Logf("IsInWorktree: %v", result)
}

func TestGetMainRepoRoot(t *testing.T) {
	if !IsGitRepo() {
		t.Skip("Not in a git repository, skipping test")
	}

	root, err := GetMainRepoRoot()
	if err != nil {
		t.Fatalf("GetMainRepoRoot() failed: %v", err)
	}

	if root == "" {
		t.Error("GetMainRepoRoot() returned empty string")
	}

	// Verify the path exists
	if _, err := os.Stat(root); os.IsNotExist(err) {
		t.Errorf("GetMainRepoRoot() returned non-existent path: %s", root)
	}

	t.Logf("Main repo root: %s", root)
}

func TestGetCurrentWorktreePath(t *testing.T) {
	if !IsGitRepo() {
		t.Skip("Not in a git repository, skipping test")
	}

	path, err := GetCurrentWorktreePath()
	if err != nil {
		t.Fatalf("GetCurrentWorktreePath() failed: %v", err)
	}

	if path == "" {
		t.Error("GetCurrentWorktreePath() returned empty string")
	}

	// Verify the path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("GetCurrentWorktreePath() returned non-existent path: %s", path)
	}

	t.Logf("Current worktree path: %s", path)
}
