package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/bshakr/ko/internal/git"
	"github.com/bshakr/ko/internal/tmux"
	"github.com/bshakr/ko/internal/validation"
	"github.com/spf13/cobra"
)

var (
	detachedCleanup bool
	detachedPath    string
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup [worktree-name]",
	Short: "Close tmux session and remove worktree",
	Long: `Closes the associated tmux window and removes the git worktree.

If no worktree name is provided and you're currently in a worktree,
it will automatically clean up the current worktree.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runCleanup,
}

func init() {
	rootCmd.AddCommand(cleanupCmd)
	cleanupCmd.Flags().BoolVar(&detachedCleanup, "detached", false, "Internal flag for detached cleanup")
	cleanupCmd.Flags().StringVar(&detachedPath, "detached-path", "", "Internal flag for worktree path")
	_ = cleanupCmd.Flags().MarkHidden("detached")
	_ = cleanupCmd.Flags().MarkHidden("detached-path")
}

func runCleanup(_ *cobra.Command, args []string) error {
	// Handle detached cleanup mode (spawned as background process)
	if detachedCleanup {
		return runDetachedCleanup()
	}

	var worktreeName string

	// If no argument provided, try to detect current worktree
	if len(args) == 0 {
		if !git.IsGitRepo() {
			return fmt.Errorf("not in a git repository\nPlease run this command from within a git repository or specify a worktree name")
		}

		if !git.IsInWorktree() {
			return fmt.Errorf("not in a worktree\nPlease specify a worktree name or run from within a worktree")
		}

		// Get the current worktree path and extract the name
		currentPath, err := git.GetCurrentWorktreePath()
		if err != nil {
			return fmt.Errorf("failed to get current worktree: %w", err)
		}

		// Extract worktree name from path (should be .ko/<name>)
		worktreeName = filepath.Base(currentPath)

		// Verify this is actually a .ko worktree by checking if parent is .ko
		parentDir := filepath.Base(filepath.Dir(currentPath))
		if parentDir != ".ko" {
			return fmt.Errorf("current worktree is not a ko worktree (not in .ko directory)\nPlease specify a worktree name explicitly")
		}

		fmt.Printf("Detected current worktree: %s\n", worktreeName)
	} else {
		worktreeName = args[0]
	}

	// Validate worktree name for security
	if err := validation.ValidateWorktreeName(worktreeName); err != nil {
		return fmt.Errorf("invalid worktree name: %w", err)
	}

	// Set up context with cancellation for long-running operations
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Signal handler goroutine - properly synchronized
	done := make(chan struct{})
	go func() {
		defer close(done)
		<-sigChan
		fmt.Println("\nOperation cancelled by user")
		cancel()
	}()
	defer func() {
		signal.Stop(sigChan)
		<-done // Wait for signal handler to finish
	}()

	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Check if we're in a git repository
	if !git.IsGitRepo() {
		return fmt.Errorf("not in a git repository\nPlease run this command from within a git repository")
	}

	// Check if worktree exists
	worktreePath := filepath.Join(currentDir, ".ko", worktreeName)
	worktreeExists := true
	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		fmt.Printf("Warning: Worktree .ko/%s not found\n", worktreeName)
		fmt.Println("Will attempt to clean up tmux window only")
		worktreeExists = false
	}

	// Check if we're running from within the worktree we're trying to cleanup
	// If so, we need to spawn a detached process to complete the cleanup after this process exits
	currentPath, err := filepath.Abs(currentDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}
	absWorktreePath, err := filepath.Abs(worktreePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute worktree path: %w", err)
	}

	// Check if current path is within the worktree being cleaned up
	isInTargetWorktree := strings.HasPrefix(currentPath, absWorktreePath+string(filepath.Separator)) ||
		currentPath == absWorktreePath

	if isInTargetWorktree && tmux.IsInTmux() {
		fmt.Println("Detected: running cleanup from within the target worktree")
		fmt.Println("Spawning detached process to complete cleanup...")
		return spawnDetachedCleanup(worktreeName, absWorktreePath)
	}

	// Close tmux window first (before removing worktree)
	tmuxWindowClosed := false
	if tmux.IsInTmux() {
		// Get repository name
		repoName, err := git.GetRepoName()
		if err != nil {
			fmt.Printf("Warning: Failed to get repository name: %v\n", err)
			repoName = ""
		}

		windowName := fmt.Sprintf("%s|%s", repoName, worktreeName)
		if err := tmux.CloseWindow(windowName, worktreeName); err != nil {
			fmt.Printf("Warning: %v\n", err)
		} else {
			fmt.Println("Tmux window closed")
			tmuxWindowClosed = true
		}
	} else {
		fmt.Println("Not in a tmux session, skipping tmux cleanup")
	}

	// Wait for shell processes to terminate after tmux window closure
	// tmux kill-window sends termination signals but doesn't wait for processes to exit
	// The shell (zsh) needs time to fully terminate and release its hold on the directory
	if tmuxWindowClosed {
		fmt.Println("Waiting for shell processes to terminate...")
		time.Sleep(1 * time.Second)
	}

	// Remove the git worktree with context (after closing tmux)
	if worktreeExists {
		fmt.Printf("Removing git worktree: .ko/%s\n", worktreeName)
		if err := git.RemoveWorktreeWithContext(ctx, worktreePath); err != nil {
			fmt.Printf("Warning: Failed to remove worktree automatically: %v\n", err)
			fmt.Printf("You may need to run: git worktree remove .ko/%s --force\n", worktreeName)
		} else {
			fmt.Println("Worktree removed successfully")
		}
	}

	fmt.Println("Cleanup complete!")
	return nil
}

// spawnDetachedCleanup spawns a detached background process to complete the cleanup
// This is necessary when cleaning up the worktree we're currently running in
func spawnDetachedCleanup(worktreeName, worktreePath string) error {
	// Get repository name for tmux window
	repoName, err := git.GetRepoName()
	if err != nil {
		fmt.Printf("Warning: Failed to get repository name: %v\n", err)
		repoName = ""
	}

	// Close the tmux window first
	windowName := fmt.Sprintf("%s|%s", repoName, worktreeName)
	if err := tmux.CloseWindow(windowName, worktreeName); err != nil {
		fmt.Printf("Warning: %v\n", err)
		// Continue anyway - try to remove worktree
	} else {
		fmt.Println("Tmux window will be closed...")
	}

	// Get the current executable path
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Spawn detached process to remove worktree after this process exits
	// Use nohup-like approach: detach from terminal, ignore signals
	// Use background context since this process should run independently
	cmd := exec.CommandContext(context.Background(), executable, "cleanup", "--detached", "--detached-path", worktreePath)

	// Detach the process completely
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,  // Create new process group
		Pgid:    0,
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to spawn detached cleanup process: %w", err)
	}

	fmt.Println("Detached cleanup process spawned successfully")
	fmt.Println("Worktree will be removed after window closure")

	// Don't wait for the command, let it run detached
	return nil
}

// runDetachedCleanup runs in the detached background process to complete cleanup
func runDetachedCleanup() error {
	if detachedPath == "" {
		return fmt.Errorf("detached-path not provided")
	}

	// Wait for parent process and shell to fully terminate
	// This needs to be longer than the normal delay since we're running detached
	time.Sleep(2 * time.Second)

	// Remove the worktree
	ctx := context.Background()
	if err := git.RemoveWorktreeWithContext(ctx, detachedPath); err != nil {
		// Log error to a file since we can't output to terminal
		logFile := filepath.Join(os.TempDir(), "ko-cleanup-error.log")
		//nolint:gosec // G306: Writing cleanup logs to temp directory is expected
		_ = os.WriteFile(logFile, []byte(fmt.Sprintf("Failed to remove worktree %s: %v\n", detachedPath, err)), 0644)
		return err
	}

	return nil
}
