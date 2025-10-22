package tmux

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bshakr/ko/internal/config"
)

// IsInTmux checks if the current session is running inside tmux
func IsInTmux() bool {
	return os.Getenv("TMUX") != ""
}

// CreateSession creates a new tmux window with 4 panes using the provided config
func CreateSession(repoName, worktreeName, worktreePath string, cfg *config.Config) error {
	return CreateSessionWithContext(context.Background(), repoName, worktreeName, worktreePath, cfg)
}

// CreateSessionWithContext creates a new tmux window with 4 panes using the provided config with cancellation support
func CreateSessionWithContext(ctx context.Context, repoName, worktreeName, worktreePath string, cfg *config.Config) error {
	if !IsInTmux() {
		return fmt.Errorf("not in a tmux session")
	}

	windowName := fmt.Sprintf("%s|%s", repoName, worktreeName)

	// Create new tmux window
	cmd := exec.CommandContext(ctx, "tmux", "new-window", "-n", windowName, "-c", worktreePath)
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.Canceled {
			return fmt.Errorf("operation cancelled")
		}
		return fmt.Errorf("failed to create tmux window: %w", err)
	}

	// Split window into 4 panes
	// First split vertically (left and right)
	if err := runTmuxCmdWithContext(ctx, "split-window", "-h", "-c", worktreePath); err != nil {
		return err
	}

	// Split left pane horizontally
	if err := runTmuxCmdWithContext(ctx, "select-pane", "-t", "0"); err != nil {
		return err
	}
	if err := runTmuxCmdWithContext(ctx, "split-window", "-v", "-c", worktreePath); err != nil {
		return err
	}

	// Split right pane horizontally
	if err := runTmuxCmdWithContext(ctx, "select-pane", "-t", "2"); err != nil {
		return err
	}
	if err := runTmuxCmdWithContext(ctx, "split-window", "-v", "-c", worktreePath); err != nil {
		return err
	}

	// Create a unique marker file path for setup completion
	setupMarker := fmt.Sprintf("/tmp/ko-setup-done-%s-%d", worktreeName, os.Getpid())

	// Run commands in each pane using config values
	// Get commands from config (use defaults if not enough specified)
	paneCommands := cfg.PaneCommands
	if len(paneCommands) < 4 {
		// Fallback to defaults if not enough commands specified
		defaults := []string{cfg.Editor, cfg.SetupScript, cfg.DevScript, "claude"}
		for i := len(paneCommands); i < 4; i++ {
			paneCommands = append(paneCommands, defaults[i])
		}
	}

	// Pane 0 (top-left): from config
	if err := sendKeysWithContext(ctx, 0, paneCommands[0]); err != nil {
		return err
	}

	// Pane 1 (bottom-left): setup script with marker
	setupCmd := fmt.Sprintf("%s && touch %s", cfg.SetupScript, setupMarker)
	if err := sendKeysWithContext(ctx, 1, setupCmd); err != nil {
		return err
	}

	// Pane 2 (top-right): wait for setup, then run dev script
	waitCmd := fmt.Sprintf("echo 'Waiting for setup to complete...'; while [ ! -f %s ]; do sleep 1; done; rm %s; %s", setupMarker, setupMarker, cfg.DevScript)
	if err := sendKeysWithContext(ctx, 2, waitCmd); err != nil {
		return err
	}

	// Pane 3 (bottom-right): from config
	if err := sendKeysWithContext(ctx, 3, paneCommands[3]); err != nil {
		return err
	}

	// Focus on the first pane
	if err := runTmuxCmdWithContext(ctx, "select-pane", "-t", "0"); err != nil {
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

// sendKeysWithContext sends keys to a specific tmux pane with cancellation support
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
