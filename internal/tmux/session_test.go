package tmux

import (
	"context"
	"testing"
	"time"

	"github.com/bshakr/koh/internal/config"
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
		SetupScript:  "",
		PaneCommands: []string{},
	}

	ctx := context.Background()
	err := CreateSessionWithContext(ctx, "test-repo", "test-worktree-0", "/tmp", cfg)
	if err != nil {
		t.Errorf("CreateSessionWithContext with no pane commands failed: %v", err)
	}

	// Cleanup
	if err := CloseWindow("test-repo", "test-worktree-0"); err != nil {
		t.Logf("Failed to close window: %v", err)
	}
}

// TestCreateSessionWithOnePaneCommand tests creating a session with setup + 1 command
func TestCreateSessionWithOnePaneCommand(t *testing.T) {
	if !IsInTmux() {
		t.Skip("Not in a tmux session, skipping test")
	}

	cfg := &config.Config{
		SetupScript:  "",
		PaneCommands: []string{"echo 'Command 1'"},
	}

	ctx := context.Background()
	err := CreateSessionWithContext(ctx, "test-repo", "test-worktree-1", "/tmp", cfg)
	if err != nil {
		t.Errorf("CreateSessionWithContext with 1 pane command failed: %v", err)
	}

	// Cleanup
	if err := CloseWindow("test-repo", "test-worktree-1"); err != nil {
		t.Logf("Failed to close window: %v", err)
	}
}

// TestCreateSessionWithTwoPaneCommands tests creating a session with setup + 2 commands
func TestCreateSessionWithTwoPaneCommands(t *testing.T) {
	if !IsInTmux() {
		t.Skip("Not in a tmux session, skipping test")
	}

	cfg := &config.Config{
		SetupScript:  "",
		PaneCommands: []string{"echo 'Command 1'", "echo 'Command 2'"},
	}

	ctx := context.Background()
	err := CreateSessionWithContext(ctx, "test-repo", "test-worktree-2", "/tmp", cfg)
	if err != nil {
		t.Errorf("CreateSessionWithContext with 2 pane commands failed: %v", err)
	}

	// Cleanup
	if err := CloseWindow("test-repo", "test-worktree-2"); err != nil {
		t.Logf("Failed to close window: %v", err)
	}
}

// TestCreateSessionWithThreePaneCommands tests creating a session with setup + 3 commands
func TestCreateSessionWithThreePaneCommands(t *testing.T) {
	if !IsInTmux() {
		t.Skip("Not in a tmux session, skipping test")
	}

	cfg := &config.Config{
		SetupScript:  "",
		PaneCommands: []string{"echo 'Command 1'", "echo 'Command 2'", "echo 'Command 3'"},
	}

	ctx := context.Background()
	err := CreateSessionWithContext(ctx, "test-repo", "test-worktree-3", "/tmp", cfg)
	if err != nil {
		t.Errorf("CreateSessionWithContext with 3 pane commands failed: %v", err)
	}

	// Cleanup
	if err := CloseWindow("test-repo", "test-worktree-3"); err != nil {
		t.Logf("Failed to close window: %v", err)
	}
}

// TestCreateSessionWithManyPaneCommands tests creating a session with setup + 5 commands
func TestCreateSessionWithManyPaneCommands(t *testing.T) {
	if !IsInTmux() {
		t.Skip("Not in a tmux session, skipping test")
	}

	cfg := &config.Config{
		SetupScript: "",
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
	if err := CloseWindow("test-repo", "test-worktree-many"); err != nil {
		t.Logf("Failed to close window: %v", err)
	}
}

// TestWindowExists tests checking if a window exists
func TestWindowExists(t *testing.T) {
	if !IsInTmux() {
		t.Skip("Not in a tmux session, skipping test")
	}

	// Test with non-existent window
	exists, err := WindowExists("nonexistent-worktree-12345")
	if err != nil {
		t.Errorf("WindowExists() error: %v", err)
	}
	if exists {
		t.Error("Expected WindowExists to return false for non-existent window")
	}
}

// TestWindowExistsWithContext tests checking if a window exists with context
func TestWindowExistsWithContext(t *testing.T) {
	if !IsInTmux() {
		t.Skip("Not in a tmux session, skipping test")
	}

	ctx := context.Background()

	// Test with non-existent window
	exists, err := WindowExistsWithContext(ctx, "nonexistent-worktree-ctx-12345")
	if err != nil {
		t.Errorf("WindowExistsWithContext() error: %v", err)
	}
	if exists {
		t.Error("Expected WindowExistsWithContext to return false for non-existent window")
	}

	// Test cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err = WindowExistsWithContext(ctx, "test-worktree")
	if err == nil {
		// Note: The command might complete before cancellation is detected
		t.Log("Command completed before cancellation (this is okay)")
	}
}

// TestSwitchToWindow tests switching to a window
func TestSwitchToWindow(t *testing.T) {
	if !IsInTmux() {
		t.Skip("Not in a tmux session, skipping test")
	}

	// Test with non-existent window should error
	err := SwitchToWindow("nonexistent-worktree-switch-12345")
	if err == nil {
		t.Error("Expected error when switching to non-existent window")
	}

	expectedErrorMsg := "no tmux window found for worktree:"
	if err != nil && !contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Expected error to contain %q, got: %v", expectedErrorMsg, err)
	}
}

// TestSwitchToWindowWithContext tests switching to a window with context
func TestSwitchToWindowWithContext(t *testing.T) {
	if !IsInTmux() {
		t.Skip("Not in a tmux session, skipping test")
	}

	ctx := context.Background()

	// Test with non-existent window should error
	err := SwitchToWindowWithContext(ctx, "nonexistent-worktree-ctx-switch-12345")
	if err == nil {
		t.Error("Expected error when switching to non-existent window")
	}
}

// TestFindWindowByWorktree tests the helper function
func TestFindWindowByWorktree(t *testing.T) {
	if !IsInTmux() {
		t.Skip("Not in a tmux session, skipping test")
	}

	ctx := context.Background()

	// Test with non-existent window
	index, name, err := findWindowByWorktree(ctx, "nonexistent-worktree-find-12345")
	if err != nil {
		t.Errorf("findWindowByWorktree() error: %v", err)
	}
	if index != "" || name != "" {
		t.Errorf("Expected empty strings for non-existent window, got index=%q, name=%q", index, name)
	}
}

// TestWindowExistsAfterCreation tests that WindowExists returns true after creating a window
func TestWindowExistsAfterCreation(t *testing.T) {
	if !IsInTmux() {
		t.Skip("Not in a tmux session, skipping test")
	}

	worktreeName := "test-exists-window"
	cfg := &config.Config{
		SetupScript:  "",
		PaneCommands: []string{},
	}

	// Create the window
	ctx := context.Background()
	err := CreateSessionWithContext(ctx, "test-repo", worktreeName, "/tmp", cfg)
	if err != nil {
		t.Fatalf("Failed to create test window: %v", err)
	}

	// Check that it exists
	exists, err := WindowExists(worktreeName)
	if err != nil {
		t.Errorf("WindowExists() error: %v", err)
	}
	if !exists {
		t.Error("Expected WindowExists to return true for created window")
	}

	// Cleanup
	if err := CloseWindow("test-repo", worktreeName); err != nil {
		t.Logf("Failed to close window: %v", err)
	}

	// Verify it no longer exists
	exists, err = WindowExists(worktreeName)
	if err != nil {
		t.Errorf("WindowExists() error after cleanup: %v", err)
	}
	if exists {
		t.Error("Expected WindowExists to return false after closing window")
	}
}

// TestGetPanesForWindow tests getting pane IDs for a window
func TestGetPanesForWindow(t *testing.T) {
	if !IsInTmux() {
		t.Skip("Not in a tmux session, skipping test")
	}

	worktreeName := "test-get-panes"
	cfg := &config.Config{
		SetupScript:  "",
		PaneCommands: []string{"echo 'pane 1'", "echo 'pane 2'"},
	}

	// Create a window with multiple panes
	ctx := context.Background()
	err := CreateSessionWithContext(ctx, "test-repo", worktreeName, "/tmp", cfg)
	if err != nil {
		t.Fatalf("Failed to create test window: %v", err)
	}

	// Get the window index
	index, _, err := findWindowByWorktree(ctx, worktreeName)
	if err != nil {
		t.Fatalf("Failed to find window: %v", err)
	}

	// Get panes for the window
	panes, err := getPanesForWindow(ctx, index)
	if err != nil {
		t.Errorf("getPanesForWindow() error: %v", err)
	}

	// We should have 3 panes (setup + 2 commands)
	expectedPanes := 3
	if len(panes) != expectedPanes {
		t.Errorf("Expected %d panes, got %d", expectedPanes, len(panes))
	}

	// Verify each pane ID starts with % (tmux pane ID format)
	for i, paneID := range panes {
		if len(paneID) == 0 || paneID[0] != '%' {
			t.Errorf("Pane %d has invalid ID format: %q", i, paneID)
		}
	}

	// Cleanup
	if err := CloseWindow("test-repo", worktreeName); err != nil {
		t.Logf("Failed to close window: %v", err)
	}
}

// TestSendCtrlCToPane tests sending Ctrl-C to a pane
func TestSendCtrlCToPane(t *testing.T) {
	if !IsInTmux() {
		t.Skip("Not in a tmux session, skipping test")
	}

	worktreeName := "test-ctrl-c"
	cfg := &config.Config{
		SetupScript:  "",
		PaneCommands: []string{},
	}

	// Create a window with one pane
	ctx := context.Background()
	err := CreateSessionWithContext(ctx, "test-repo", worktreeName, "/tmp", cfg)
	if err != nil {
		t.Fatalf("Failed to create test window: %v", err)
	}

	// Get the window index and panes
	index, _, err := findWindowByWorktree(ctx, worktreeName)
	if err != nil {
		t.Fatalf("Failed to find window: %v", err)
	}

	panes, err := getPanesForWindow(ctx, index)
	if err != nil {
		t.Fatalf("Failed to get panes: %v", err)
	}

	if len(panes) == 0 {
		t.Fatal("Expected at least one pane")
	}

	// Send Ctrl-C to the first pane
	err = sendCtrlCToPane(ctx, panes[0])
	if err != nil {
		t.Errorf("sendCtrlCToPane() error: %v", err)
	}

	// Cleanup
	if err := CloseWindow("test-repo", worktreeName); err != nil {
		t.Logf("Failed to close window: %v", err)
	}
}

// TestGetPanesForWindowNonExistent tests error handling for non-existent window
func TestGetPanesForWindowNonExistent(t *testing.T) {
	if !IsInTmux() {
		t.Skip("Not in a tmux session, skipping test")
	}

	ctx := context.Background()
	// Use a very high window index that should not exist
	panes, err := getPanesForWindow(ctx, "999999")
	if err == nil {
		t.Error("Expected error for non-existent window, got nil")
	}
	if len(panes) != 0 {
		t.Errorf("Expected 0 panes for non-existent window, got %d", len(panes))
	}
}

// TestCloseWindowWithCtrlC tests that CloseWindow sends Ctrl-C before killing
func TestCloseWindowWithCtrlC(t *testing.T) {
	if !IsInTmux() {
		t.Skip("Not in a tmux session, skipping test")
	}

	worktreeName := "test-close-with-ctrl-c"
	cfg := &config.Config{
		SetupScript:  "",
		PaneCommands: []string{"sleep 10", "sleep 20"},
	}

	// Create a window with panes running sleep commands
	ctx := context.Background()
	err := CreateSessionWithContext(ctx, "test-repo", worktreeName, "/tmp", cfg)
	if err != nil {
		t.Fatalf("Failed to create test window: %v", err)
	}

	// Verify window exists
	exists, err := WindowExists(worktreeName)
	if err != nil {
		t.Fatalf("WindowExists() error: %v", err)
	}
	if !exists {
		t.Fatal("Window should exist after creation")
	}

	// Close the window (should send Ctrl-C to all panes before killing)
	err = CloseWindow("test-repo", worktreeName)
	if err != nil {
		t.Errorf("CloseWindow() error: %v", err)
	}

	// Verify window no longer exists
	exists, err = WindowExists(worktreeName)
	if err != nil {
		t.Errorf("WindowExists() error after close: %v", err)
	}
	if exists {
		t.Error("Window should not exist after closing")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
