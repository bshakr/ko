package cmd

import (
	"testing"

	"github.com/bshakr/ko/internal/validation"
)

func TestWorktreeNameValidationInNew(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name:      "valid worktree name",
			input:     "feature-branch",
			shouldErr: false,
		},
		{
			name:      "path traversal attempt",
			input:     "../../../etc/passwd",
			shouldErr: true,
		},
		{
			name:      "dot dot",
			input:     "..",
			shouldErr: true,
		},
		{
			name:      "slash in name",
			input:     "feature/branch",
			shouldErr: true,
		},
		{
			name:      "backslash in name",
			input:     "feature\\branch",
			shouldErr: true,
		},
		{
			name:      "hidden traversal",
			input:     "feature..hack",
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

// TestNewCommandStructure verifies the command is properly configured
func TestNewCommandStructure(t *testing.T) {
	if newCmd == nil {
		t.Fatal("newCmd is nil")
	}

	if newCmd.Use != "new <worktree-name>" {
		t.Errorf("Expected Use 'new <worktree-name>', got %q", newCmd.Use)
	}

	if newCmd.Short == "" {
		t.Error("Short description is empty")
	}

	if newCmd.Long == "" {
		t.Error("Long description is empty")
	}

	if newCmd.RunE == nil {
		t.Error("RunE is nil")
	}
}
