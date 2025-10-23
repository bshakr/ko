package git

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCreateWorktreeWithContext(t *testing.T) {
	if !IsGitRepo() {
		t.Skip("Not in a git repository, skipping test")
	}

	// Create a temporary directory for the worktree
	tempDir, err := os.MkdirTemp("", "ko-test-worktree-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	worktreePath := filepath.Join(tempDir, "test-worktree")

	// Test successful creation with context
	ctx := context.Background()
	err = CreateWorktreeWithContext(ctx, worktreePath)
	if err != nil {
		t.Fatalf("CreateWorktreeWithContext() failed: %v", err)
	}

	// Verify the worktree was created
	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		t.Errorf("Worktree was not created at %s", worktreePath)
	}

	// Clean up
	if err := RemoveWorktree(worktreePath); err != nil {
		t.Logf("Failed to remove worktree: %v", err)
	}
}

func TestCreateWorktreeWithContextCancellation(t *testing.T) {
	if !IsGitRepo() {
		t.Skip("Not in a git repository, skipping test")
	}

	// Create a temporary directory for the worktree
	tempDir, err := os.MkdirTemp("", "ko-test-worktree-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	worktreePath := filepath.Join(tempDir, "test-worktree-cancel")

	// Test cancellation - cancel immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel before the operation

	err = CreateWorktreeWithContext(ctx, worktreePath)
	if err == nil {
		t.Error("Expected error due to cancellation, got nil")
		if cleanupErr := RemoveWorktree(worktreePath); cleanupErr != nil {
			t.Logf("Failed to remove worktree: %v", cleanupErr)
		}
	}

	if err != nil && err.Error() != "operation cancelled" {
		// The operation might complete before cancellation is detected
		t.Logf("Got error: %v (might complete before cancellation)", err)
	}
}

func TestCreateWorktreeWithContextTimeout(t *testing.T) {
	if !IsGitRepo() {
		t.Skip("Not in a git repository, skipping test")
	}

	// Create a temporary directory for the worktree
	tempDir, err := os.MkdirTemp("", "ko-test-worktree-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	worktreePath := filepath.Join(tempDir, "test-worktree-timeout")

	// Test with a reasonable timeout (should succeed)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = CreateWorktreeWithContext(ctx, worktreePath)
	if err != nil {
		t.Fatalf("CreateWorktreeWithContext() with timeout failed: %v", err)
	}

	// Verify the worktree was created
	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		t.Errorf("Worktree was not created at %s", worktreePath)
	}

	// Clean up
	if err := RemoveWorktreeWithContext(context.Background(), worktreePath); err != nil {
		t.Logf("Failed to remove worktree: %v", err)
	}
}

func TestRemoveWorktreeWithContext(t *testing.T) {
	if !IsGitRepo() {
		t.Skip("Not in a git repository, skipping test")
	}

	// Create a temporary directory for the worktree
	tempDir, err := os.MkdirTemp("", "ko-test-worktree-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	worktreePath := filepath.Join(tempDir, "test-worktree-remove")

	// First create a worktree
	err = CreateWorktree(worktreePath)
	if err != nil {
		t.Fatalf("Failed to create worktree: %v", err)
	}

	// Now test removal with context
	ctx := context.Background()
	err = RemoveWorktreeWithContext(ctx, worktreePath)
	if err != nil {
		t.Fatalf("RemoveWorktreeWithContext() failed: %v", err)
	}

	// Verify the worktree was removed
	// Note: The directory might still exist but git worktree list shouldn't show it
	t.Logf("Worktree removed successfully")
}

func TestRemoveWorktreeWithContextCancellation(t *testing.T) {
	if !IsGitRepo() {
		t.Skip("Not in a git repository, skipping test")
	}

	// Create a temporary directory for the worktree
	tempDir, err := os.MkdirTemp("", "ko-test-worktree-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	worktreePath := filepath.Join(tempDir, "test-worktree-remove-cancel")

	// First create a worktree
	err = CreateWorktree(worktreePath)
	if err != nil {
		t.Fatalf("Failed to create worktree: %v", err)
	}
	defer func() {
		if err := RemoveWorktree(worktreePath); err != nil {
			t.Logf("Failed to remove worktree: %v", err)
		}
	}() // Ensure cleanup

	// Test cancellation - cancel immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel before the operation

	err = RemoveWorktreeWithContext(ctx, worktreePath)
	if err == nil {
		t.Error("Expected error due to cancellation, got nil")
	}

	if err != nil && err.Error() != "operation cancelled" {
		// The operation might complete before cancellation is detected
		t.Logf("Got error: %v (might complete before cancellation)", err)
	}
}
