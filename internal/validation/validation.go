// Package validation provides security-focused input validation for ko.
//
// This package validates user-supplied input to prevent security issues such as:
//   - Path traversal attacks (using .. or absolute paths)
//   - Special characters that could cause issues in shell commands
//   - Reserved system names that could cause conflicts
//   - Overly long input that could cause buffer issues
//
// All user-supplied worktree names must pass through ValidateWorktreeName
// before being used in file operations or shell commands.
//
// The validation is designed to be strict and cross-platform compatible,
// rejecting potentially dangerous input even on systems where it might be safe.
package validation

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ValidateWorktreeName validates that a worktree name is safe to use
func ValidateWorktreeName(name string) error {
	if name == "" {
		return fmt.Errorf("worktree name cannot be empty")
	}

	// Check length (reasonable limit)
	if len(name) > 255 {
		return fmt.Errorf("worktree name too long (max 255 characters)")
	}

	// Check for path separators
	if strings.ContainsAny(name, "/\\") {
		return fmt.Errorf("worktree name cannot contain path separators (/ or \\)")
	}

	// Check for path traversal attempts
	if name == "." || name == ".." {
		return fmt.Errorf("worktree name cannot be '.' or '..'")
	}

	if strings.Contains(name, "..") {
		return fmt.Errorf("worktree name cannot contain '..'")
	}

	// Check for special characters that could cause issues
	if strings.ContainsAny(name, "\x00\n\r\t") {
		return fmt.Errorf("worktree name contains invalid characters")
	}

	// Ensure it's not trying to escape using filepath operations
	cleaned := filepath.Clean(name)
	if cleaned != name {
		return fmt.Errorf("worktree name contains invalid path components")
	}

	// Check for reserved names on Windows (even if we're not on Windows, be safe)
	reserved := []string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4",
		"COM5", "COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3",
		"LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}
	upperName := strings.ToUpper(name)
	for _, r := range reserved {
		if upperName == r {
			return fmt.Errorf("worktree name cannot be a reserved system name")
		}
	}

	return nil
}
