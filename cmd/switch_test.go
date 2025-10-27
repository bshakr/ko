package cmd

import (
	"testing"

	"github.com/bshakr/ko/internal/validation"
)

func TestWorktreeNameValidationInSwitch(t *testing.T) {
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
			name:      "valid with numbers",
			input:     "feature-123",
			shouldErr: false,
		},
		{
			name:      "valid with underscores",
			input:     "feature_branch",
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
			name:      "single dot",
			input:     ".",
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
		{
			name:      "null byte",
			input:     "feature\x00branch",
			shouldErr: true,
		},
		{
			name:      "newline",
			input:     "feature\nbranch",
			shouldErr: true,
		},
		{
			name:      "empty string",
			input:     "",
			shouldErr: true,
		},
		{
			name:      "reserved name CON",
			input:     "CON",
			shouldErr: true,
		},
		{
			name:      "reserved name NUL",
			input:     "NUL",
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
				t.Errorf("Expected no error for input %q, got: %v", tt.input, err)
			}
		})
	}
}
