package validation

import (
	"strings"
	"testing"
)

func TestValidateWorktreeName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		errMsg    string
		shouldErr bool
	}{
		{
			name:      "valid simple name",
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
			name:      "empty name",
			input:     "",
			shouldErr: true,
			errMsg:    "cannot be empty",
		},
		{
			name:      "path traversal with ../",
			input:     "../etc/passwd",
			shouldErr: true,
			errMsg:    "path separators",
		},
		{
			name:      "path traversal with ..",
			input:     "..",
			shouldErr: true,
			errMsg:    "cannot be '.' or '..'",
		},
		{
			name:      "path traversal hidden",
			input:     "feature..branch",
			shouldErr: true,
			errMsg:    "cannot contain '..'",
		},
		{
			name:      "forward slash",
			input:     "feature/branch",
			shouldErr: true,
			errMsg:    "path separators",
		},
		{
			name:      "backslash",
			input:     "feature\\branch",
			shouldErr: true,
			errMsg:    "path separators",
		},
		{
			name:      "null byte",
			input:     "feature\x00branch",
			shouldErr: true,
			errMsg:    "invalid characters",
		},
		{
			name:      "newline",
			input:     "feature\nbranch",
			shouldErr: true,
			errMsg:    "invalid characters",
		},
		{
			name:      "too long name",
			input:     strings.Repeat("a", 256),
			shouldErr: true,
			errMsg:    "too long",
		},
		{
			name:      "reserved name CON",
			input:     "CON",
			shouldErr: true,
			errMsg:    "reserved system name",
		},
		{
			name:      "reserved name aux",
			input:     "aux",
			shouldErr: true,
			errMsg:    "reserved system name",
		},
		{
			name:      "just dot",
			input:     ".",
			shouldErr: true,
			errMsg:    "cannot be '.' or '..'",
		},
		{
			name:      "valid with hyphens",
			input:     "my-feature-branch",
			shouldErr: false,
		},
		{
			name:      "valid with mixed case",
			input:     "MyFeatureBranch",
			shouldErr: false,
		},
		{
			name:      "tab character",
			input:     "feature\tbranch",
			shouldErr: true,
			errMsg:    "invalid characters",
		},
		{
			name:      "carriage return",
			input:     "feature\rbranch",
			shouldErr: true,
			errMsg:    "invalid characters",
		},
		{
			name:      "exactly 255 chars (boundary)",
			input:     strings.Repeat("a", 255),
			shouldErr: false,
		},
		{
			name:      "reserved name lowercase prn",
			input:     "prn",
			shouldErr: true,
			errMsg:    "reserved system name",
		},
		{
			name:      "valid name with numbers at start",
			input:     "123-feature",
			shouldErr: false,
		},
		{
			name:      "path with absolute path attempt",
			input:     "/etc/passwd",
			shouldErr: true,
			errMsg:    "path separators",
		},
		{
			name:      "windows path",
			input:     "C:\\Windows\\System32",
			shouldErr: true,
			errMsg:    "path separators",
		},
		{
			name:      "hidden path traversal in middle",
			input:     "a/../etc",
			shouldErr: true,
			errMsg:    "path separators",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateWorktreeName(tt.input)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errMsg)
				} else if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}
