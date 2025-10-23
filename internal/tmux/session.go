// Package tmux provides functionality for managing tmux sessions and windows.
//
// # Security Note: Command Injection Trust Model
//
// This package executes commands from the user's .koconfig file without sanitization.
// This is intentional and follows these security principles:
//
//  1. Configuration is Local and User-Controlled:
//     - The .koconfig file is stored in the user's repository
//     - Users explicitly create and edit this file
//     - This is equivalent to running shell scripts in the repository
//
//  2. No Remote/Untrusted Input:
//     - Commands come from local configuration files only
//     - No network inputs or user-supplied runtime arguments are executed
//     - Worktree names are validated separately (see validation package)
//
//  3. Trust Model is Consistent with Git:
//     - Similar to .git/hooks, which execute arbitrary scripts
//     - Users already trust their repository contents
//     - If the repository is compromised, the system is already at risk
//
//  4. Protection at the Configuration Layer:
//     - The 'ko init' command shows users exactly what will be executed
//     - Configuration is human-readable JSON
//     - Users should review .koconfig before committing
//
// WARNING: Do not extend this package to execute commands from:
// - Network sources
// - User input at runtime (beyond validated worktree names)
// - Untrusted or external configuration sources
//
// If this trust model is unacceptable for your use case, review and modify
// the .koconfig file before using ko commands.
package tmux

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bshakr/ko/internal/config"
	"github.com/bshakr/ko/internal/git"
)

// IsInTmux checks if the current session is running inside tmux
func IsInTmux() bool {
	return os.Getenv("TMUX") != ""
}

// ensureSetupScript checks if the setup script exists in the worktree.
// If not, it looks for it in the main repo root and copies it to the worktree.
// Returns an error if the script cannot be found or copied.
func ensureSetupScript(worktreePath, setupScript string) error {
	// If setup script is empty, nothing to do
	if setupScript == "" {
		return nil
	}

	// Check if the setup script path is absolute
	var scriptPath string
	if filepath.IsAbs(setupScript) {
		scriptPath = setupScript
	} else {
		scriptPath = filepath.Join(worktreePath, setupScript)
	}

	// Check if the script exists in the worktree
	if _, err := os.Stat(scriptPath); err == nil {
		// Script exists in worktree, nothing to do
		return nil
	}

	// Script doesn't exist in worktree, try to copy from main repo
	// Get the main repo root
	mainRepoRoot, err := git.GetMainRepoRoot()
	if err != nil {
		return fmt.Errorf("failed to get main repo root: %w", err)
	}

	// Check if the script exists in the main repo root
	mainRepoScriptPath := filepath.Join(mainRepoRoot, setupScript)
	if _, err := os.Stat(mainRepoScriptPath); os.IsNotExist(err) {
		// Script doesn't exist in main repo either
		return fmt.Errorf("setup script not found in worktree or main repo: %s", setupScript)
	}

	// Copy the script from main repo to worktree
	if err := copyFile(mainRepoScriptPath, scriptPath); err != nil {
		return fmt.Errorf("failed to copy setup script from main repo: %w", err)
	}

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	// Open source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer func() {
		_ = sourceFile.Close() // Ignore error in defer
	}()

	// Get source file info to preserve permissions
	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat source file: %w", err)
	}

	// Create destination directory if it doesn't exist
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Create destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func() {
		_ = destFile.Close() // Ignore error in defer
	}()

	// Copy the file content
	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	// Preserve file permissions
	if err := os.Chmod(dst, sourceInfo.Mode()); err != nil {
		return fmt.Errorf("failed to set file permissions: %w", err)
	}

	return nil
}

// CreateSession creates a new tmux window with dynamically created panes based on the provided config
func CreateSession(repoName, worktreeName, worktreePath string, cfg *config.Config) error {
	return CreateSessionWithContext(context.Background(), repoName, worktreeName, worktreePath, cfg)
}

// CreateSessionWithContext creates a new tmux window with dynamically created panes based on config
func CreateSessionWithContext(ctx context.Context, repoName, worktreeName, worktreePath string, cfg *config.Config) error {
	if !IsInTmux() {
		return fmt.Errorf("not in a tmux session")
	}

	// Ensure the setup script is available (copy from main repo if needed)
	if err := ensureSetupScript(worktreePath, cfg.SetupScript); err != nil {
		return fmt.Errorf("failed to ensure setup script: %w", err)
	}

	// Get the pane base index from tmux configuration
	paneBaseIndex, err := getPaneBaseIndex(ctx)
	if err != nil {
		return fmt.Errorf("failed to get pane base index: %w", err)
	}

	windowName := fmt.Sprintf("%s|%s", repoName, worktreeName)

	// Create new tmux window with setup script
	cmd := exec.CommandContext(ctx, "tmux", "new-window", "-n", windowName, "-c", worktreePath)
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.Canceled {
			return fmt.Errorf("operation cancelled")
		}
		return fmt.Errorf("failed to create tmux window: %w", err)
	}

	// Create panes dynamically based on pane_commands
	// Layout strategy:
	// - Pane 0 (baseIndex): Setup (always)
	// - Pane 1 (baseIndex+1): First command - side by side with setup (vertical split)
	// - Pane 2 (baseIndex+2): Second command - under setup (split pane 0 horizontally)
	// - Pane 3 (baseIndex+3): Third command - under first command (split pane 1 horizontally)
	// - Pane 4 (baseIndex+4): Fourth command - under second command (split pane 2 horizontally)
	// - Continue pattern: each new pane splits the pane created 2 steps before
	numPaneCommands := len(cfg.PaneCommands)

	// If there are pane commands, create additional panes
	if numPaneCommands > 0 {
		// First pane command: split vertically to create side-by-side layout (setup | command1)
		if err := runTmuxCmdWithContext(ctx, "split-window", "-h", "-c", worktreePath); err != nil {
			return err
		}

		// Additional pane commands: split existing panes horizontally
		// Pattern: split pane (i-1) to create pane (i+1)
		for i := 1; i < numPaneCommands; i++ {
			// For second command (i=1): split pane 0 (setup)
			// For third command (i=2): split pane 1 (first command)
			// For fourth command (i=3): split pane 2 (second command)
			// General formula: target pane index = paneBaseIndex + (i - 1)
			targetPaneIdx := paneBaseIndex + (i - 1)
			targetPane := fmt.Sprintf("%d", targetPaneIdx)

			// Select the target pane and split horizontally
			if err := runTmuxCmdWithContext(ctx, "select-pane", "-t", targetPane); err != nil {
				return err
			}
			if err := runTmuxCmdWithContext(ctx, "split-window", "-v", "-c", worktreePath); err != nil {
				return err
			}
		}
	}

	// Send commands to panes
	// Pane 0: Setup script (always)
	if cfg.SetupScript != "" {
		if err := sendKeysWithContext(ctx, paneBaseIndex, cfg.SetupScript); err != nil {
			return err
		}
	}

	// Panes 1+: Pane commands
	// The pane mapping is:
	// - cfg.PaneCommands[0] -> pane baseIndex+1
	// - cfg.PaneCommands[1] -> pane baseIndex+2
	// - cfg.PaneCommands[n] -> pane baseIndex+n+1
	for i, cmd := range cfg.PaneCommands {
		paneIdx := paneBaseIndex + i + 1
		if err := sendKeysWithContext(ctx, paneIdx, cmd); err != nil {
			return err
		}
	}

	// Focus on the first pane (setup)
	firstPane := fmt.Sprintf("%d", paneBaseIndex)
	if err := runTmuxCmdWithContext(ctx, "select-pane", "-t", firstPane); err != nil {
		return err
	}

	return nil
}

// CloseWindow closes a tmux window by name
func CloseWindow(windowName, worktreeName string) error {
	// Find the window index with the worktree name
	cmd := exec.Command("tmux", "list-windows", "-F", "#{window_index}:#{window_name}")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to list tmux windows: %w", err)
	}

	var windowIndex string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasSuffix(line, "|"+worktreeName) {
			parts := strings.Split(line, ":")
			if len(parts) > 0 {
				windowIndex = parts[0]
				break
			}
		}
	}

	if windowIndex == "" {
		return fmt.Errorf("no tmux window found with name: %s", worktreeName)
	}

	// Kill the window
	cmd = exec.Command("tmux", "kill-window", "-t", windowIndex)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to close tmux window: %w", err)
	}

	return nil
}

// runTmuxCmd runs a tmux command with the given arguments
func runTmuxCmd(args ...string) error {
	return runTmuxCmdWithContext(context.Background(), args...)
}

// runTmuxCmdWithContext runs a tmux command with the given arguments with cancellation support
func runTmuxCmdWithContext(ctx context.Context, args ...string) error {
	cmd := exec.CommandContext(ctx, "tmux", args...)
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.Canceled {
			return fmt.Errorf("operation cancelled")
		}
		return fmt.Errorf("tmux command failed (%v): %w", args, err)
	}
	return nil
}

// sendKeys sends keys to a specific tmux pane
func sendKeys(pane int, keys string) error {
	return sendKeysWithContext(context.Background(), pane, keys)
}

// sendKeysWithContext sends keys to a specific tmux pane with cancellation support.
//
// Security: The 'keys' parameter contains commands from the user's .koconfig file.
// These commands are trusted local configuration (see package documentation for
// the security trust model). The keys are passed to tmux which executes them
// in the shell context of the pane.
func sendKeysWithContext(ctx context.Context, pane int, keys string) error {
	paneTarget := fmt.Sprintf("%d", pane)
	cmd := exec.CommandContext(ctx, "tmux", "send-keys", "-t", paneTarget, keys, "C-m")
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.Canceled {
			return fmt.Errorf("operation cancelled")
		}
		return fmt.Errorf("failed to send keys to pane %d: %w", pane, err)
	}
	return nil
}

// getPaneBaseIndex retrieves the pane-base-index setting from tmux configuration
func getPaneBaseIndex(ctx context.Context) (int, error) {
	cmd := exec.CommandContext(ctx, "tmux", "show-options", "-gv", "pane-base-index")
	output, err := cmd.Output()
	if err != nil {
		// If the option is not set, default to 0
		return 0, nil
	}

	var baseIndex int
	_, err = fmt.Sscanf(strings.TrimSpace(string(output)), "%d", &baseIndex)
	if err != nil {
		return 0, fmt.Errorf("failed to parse pane-base-index: %w", err)
	}

	return baseIndex, nil
}
