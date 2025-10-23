package tmux

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bshakr/ko/internal/config"
	"github.com/bshakr/ko/internal/git"
)

func TestIsInTmux(t *testing.T) {
	// This test just verifies the function runs without error
	result := IsInTmux()
	t.Logf("IsInTmux: %v", result)

	// Note: This will return false when running tests outside tmux
	// and true when running inside tmux
}

func TestRunTmuxCmdWithContext(t *testing.T) {
	if !IsInTmux() {
		t.Skip("Not in a tmux session, skipping test")
	}

	// Test a simple tmux command with context
	ctx := context.Background()
	err := runTmuxCmdWithContext(ctx, "display-message", "-p", "test")
	if err != nil {
		t.Errorf("runTmuxCmdWithContext() failed: %v", err)
	}
}

func TestRunTmuxCmdWithContextCancellation(t *testing.T) {
	if !IsInTmux() {
		t.Skip("Not in a tmux session, skipping test")
	}

	// Test cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := runTmuxCmdWithContext(ctx, "display-message", "-p", "test")
	if err == nil {
		t.Error("Expected error due to cancellation, got nil")
	}

	if err != nil && err.Error() != "operation cancelled" {
		// The operation might complete before cancellation is detected
		t.Logf("Got error: %v (might complete before cancellation)", err)
	}
}

func TestRunTmuxCmdWithContextTimeout(t *testing.T) {
	if !IsInTmux() {
		t.Skip("Not in a tmux session, skipping test")
	}

	// Test with a reasonable timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := runTmuxCmdWithContext(ctx, "display-message", "-p", "test")
	if err != nil {
		t.Errorf("runTmuxCmdWithContext() with timeout failed: %v", err)
	}
}

func TestSendKeysWithContext(t *testing.T) {
	if !IsInTmux() {
		t.Skip("Not in a tmux session, skipping test")
	}

	// We can't easily test this without creating actual panes
	// So we'll just test that the function exists and has the right signature
	ctx := context.Background()

	// Try to send keys to pane 0 (should fail if we're not in the right window)
	err := sendKeysWithContext(ctx, 0, "echo test")
	// We expect this might fail depending on the tmux setup
	t.Logf("sendKeysWithContext result: %v", err)
}

func TestSendKeysWithContextCancellation(t *testing.T) {
	if !IsInTmux() {
		t.Skip("Not in a tmux session, skipping test")
	}

	// Test cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := sendKeysWithContext(ctx, 0, "echo test")
	if err == nil {
		t.Error("Expected error due to cancellation, got nil")
	}

	if err != nil && err.Error() != "operation cancelled" {
		// The operation might complete or fail for other reasons
		t.Logf("Got error: %v (might complete before cancellation or fail for other reasons)", err)
	}
}

func TestCloseWindow(t *testing.T) {
	if !IsInTmux() {
		t.Skip("Not in a tmux session, skipping test")
	}

	// We can't easily test this without creating actual windows
	// Just verify the function exists
	err := CloseWindow("test-window", "test-worktree")

	// We expect this to fail since the window doesn't exist
	if err == nil {
		t.Error("Expected error for non-existent window, got nil")
	}

	t.Logf("CloseWindow error (expected): %v", err)
}

// Test that the backwards-compatible functions still work
func TestBackwardsCompatibility(t *testing.T) {
	if !IsInTmux() {
		t.Skip("Not in a tmux session, skipping test")
	}

	// Test that the non-context versions still work by delegating to context versions
	err := runTmuxCmd("display-message", "-p", "test")
	if err != nil {
		t.Errorf("runTmuxCmd() (backwards compatible) failed: %v", err)
	}

	// Test sendKeys backwards compatibility
	err = sendKeys(0, "echo test")
	// Might fail depending on pane setup, but should not panic
	t.Logf("sendKeys result: %v", err)
}

// TestCreateSessionWithNoPaneCommands tests creating a session with only setup script
func TestCreateSessionWithNoPaneCommands(t *testing.T) {
	if !IsInTmux() {
		t.Skip("Not in a tmux session, skipping test")
	}

	cfg := &config.Config{
		SetupScript:  "echo 'Setup'",
		PaneCommands: []string{},
	}

	ctx := context.Background()
	err := CreateSessionWithContext(ctx, "test-repo", "test-worktree-0", "/tmp", cfg)
	if err != nil {
		t.Errorf("CreateSessionWithContext with no pane commands failed: %v", err)
	}

	// Cleanup
	CloseWindow("test-repo", "test-worktree-0")
}

// TestCreateSessionWithOnePaneCommand tests creating a session with setup + 1 command
func TestCreateSessionWithOnePaneCommand(t *testing.T) {
	if !IsInTmux() {
		t.Skip("Not in a tmux session, skipping test")
	}

	cfg := &config.Config{
		SetupScript:  "echo 'Setup'",
		PaneCommands: []string{"echo 'Command 1'"},
	}

	ctx := context.Background()
	err := CreateSessionWithContext(ctx, "test-repo", "test-worktree-1", "/tmp", cfg)
	if err != nil {
		t.Errorf("CreateSessionWithContext with 1 pane command failed: %v", err)
	}

	// Cleanup
	CloseWindow("test-repo", "test-worktree-1")
}

// TestCreateSessionWithTwoPaneCommands tests creating a session with setup + 2 commands
func TestCreateSessionWithTwoPaneCommands(t *testing.T) {
	if !IsInTmux() {
		t.Skip("Not in a tmux session, skipping test")
	}

	cfg := &config.Config{
		SetupScript:  "echo 'Setup'",
		PaneCommands: []string{"echo 'Command 1'", "echo 'Command 2'"},
	}

	ctx := context.Background()
	err := CreateSessionWithContext(ctx, "test-repo", "test-worktree-2", "/tmp", cfg)
	if err != nil {
		t.Errorf("CreateSessionWithContext with 2 pane commands failed: %v", err)
	}

	// Cleanup
	CloseWindow("test-repo", "test-worktree-2")
}

// TestCreateSessionWithThreePaneCommands tests creating a session with setup + 3 commands
func TestCreateSessionWithThreePaneCommands(t *testing.T) {
	if !IsInTmux() {
		t.Skip("Not in a tmux session, skipping test")
	}

	cfg := &config.Config{
		SetupScript:  "echo 'Setup'",
		PaneCommands: []string{"echo 'Command 1'", "echo 'Command 2'", "echo 'Command 3'"},
	}

	ctx := context.Background()
	err := CreateSessionWithContext(ctx, "test-repo", "test-worktree-3", "/tmp", cfg)
	if err != nil {
		t.Errorf("CreateSessionWithContext with 3 pane commands failed: %v", err)
	}

	// Cleanup
	CloseWindow("test-repo", "test-worktree-3")
}

// TestCreateSessionWithManyPaneCommands tests creating a session with setup + 5 commands
func TestCreateSessionWithManyPaneCommands(t *testing.T) {
	if !IsInTmux() {
		t.Skip("Not in a tmux session, skipping test")
	}

	cfg := &config.Config{
		SetupScript: "echo 'Setup'",
		PaneCommands: []string{
			"echo 'Command 1'",
			"echo 'Command 2'",
			"echo 'Command 3'",
			"echo 'Command 4'",
			"echo 'Command 5'",
		},
	}

	ctx := context.Background()
	err := CreateSessionWithContext(ctx, "test-repo", "test-worktree-many", "/tmp", cfg)
	if err != nil {
		t.Errorf("CreateSessionWithContext with many pane commands failed: %v", err)
	}

	// Cleanup
	CloseWindow("test-repo", "test-worktree-many")
}
