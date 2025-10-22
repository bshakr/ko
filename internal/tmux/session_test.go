package tmux

import (
	"context"
	"testing"
	"time"
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
