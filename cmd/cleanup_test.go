package cmd

import (
	"testing"

	"github.com/bshakr/ko/internal/validation"
)

func TestWorktreeNameValidationInCleanup(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name:      "valid worktree name",
			input:     "old-feature",
			shouldErr: false,
		},
		{
			name:      "path traversal attempt",
			input:     "../../../tmp",
			shouldErr: true,
		},
		{
			name:      "dot dot",
			input:     "..",
			shouldErr: true,
		},
		{
			name:      "absolute path",
			input:     "/tmp/worktree",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateWorktreeName(tt.input)
			if tt.shouldErr && err == nil {
				t.Errorf("Expected error for input %q, got nil", tt.input)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Expected no error for input %q, got %v", tt.input, err)
			}
		})
	}
}

// TestCleanupCommandStructure verifies the command is properly configured
func TestCleanupCommandStructure(t *testing.T) {
	if cleanupCmd == nil {
		t.Fatal("cleanupCmd is nil")
	}

	if cleanupCmd.Use != "cleanup [worktree-name]" {
		t.Errorf("Expected Use 'cleanup [worktree-name]', got %q", cleanupCmd.Use)
	}

	if cleanupCmd.Short == "" {
		t.Error("Short description is empty")
	}

	if cleanupCmd.Long == "" {
		t.Error("Long description is empty")
	}

	if cleanupCmd.RunE == nil {
		t.Error("RunE is nil")
	}
}

// TestCleanupAutoDetection tests that cleanup can auto-detect current worktree
func TestCleanupAutoDetection(t *testing.T) {
	// This is more of a documentation test showing the expected behavior
	// Actual behavior testing requires being in a worktree which we can't
	// guarantee in unit tests

	// When no args provided and in a worktree -> should detect and cleanup
	// When no args provided and NOT in a worktree -> should error
	// When arg provided -> should use provided name

	t.Log("Cleanup command supports auto-detection of current worktree")
	t.Log("Run 'ko cleanup' from within a worktree to test this feature")
}
