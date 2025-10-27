package cmd

import (
	"runtime"
	"testing"
	"time"

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

// TestRunNewFailsGracefullyWithoutHanging verifies that runNew returns an error
// without hanging when preconditions aren't met (no git repo, no config, etc.)
func TestRunNewFailsGracefullyWithoutHanging(t *testing.T) {
	tests := []struct {
		name         string
		worktreeName string
		wantTimeout  bool
	}{
		{
			name:         "invalid worktree name fails quickly",
			worktreeName: "../../../etc/passwd",
			wantTimeout:  false,
		},
		{
			name:         "path traversal fails quickly",
			worktreeName: "..",
			wantTimeout:  false,
		},
		{
			name:         "path separator fails quickly",
			worktreeName: "feature/branch",
			wantTimeout:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Track initial goroutine count
			initialGoroutines := runtime.NumGoroutine()

			// Run command in goroutine with timeout
			done := make(chan error, 1)
			go func() {
				args := []string{tt.worktreeName}
				done <- runNew(newCmd, args)
			}()

			// Wait for command to complete or timeout
			select {
			case err := <-done:
				if err == nil {
					t.Error("Expected error due to invalid input or missing preconditions, got nil")
				}
				// Expected: command failed with error
				t.Logf("Command failed as expected: %v", err)
			case <-time.After(3 * time.Second):
				if tt.wantTimeout {
					t.Log("Command timed out as expected")
				} else {
					t.Fatal("Command hung and did not complete within 3 seconds")
				}
			}

			// Give goroutines time to clean up
			time.Sleep(100 * time.Millisecond)

			// Check for goroutine leaks
			finalGoroutines := runtime.NumGoroutine()
			if finalGoroutines > initialGoroutines+2 {
				t.Errorf("Potential goroutine leak: initial=%d, final=%d", initialGoroutines, finalGoroutines)
			}
		})
	}
}

// TestRunNewValidatesInputBeforeSettingUpSignals verifies that input validation
// happens before signal handling setup, ensuring fast failure
func TestRunNewValidatesInputBeforeSettingUpSignals(t *testing.T) {
	invalidNames := []string{
		"../../../etc/passwd",
		"..",
		"feature/branch",
		"feature\\branch",
	}

	for _, name := range invalidNames {
		t.Run(name, func(t *testing.T) {
			start := time.Now()

			args := []string{name}
			err := runNew(newCmd, args)

			duration := time.Since(start)

			// Should fail quickly (within 100ms for validation)
			if duration > 100*time.Millisecond {
				t.Errorf("Validation took too long: %v", duration)
			}

			// Should return error
			if err == nil {
				t.Error("Expected validation error, got nil")
			}
		})
	}
}
