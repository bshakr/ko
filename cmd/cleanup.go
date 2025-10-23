package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/bshakr/ko/internal/git"
	"github.com/bshakr/ko/internal/tmux"
	"github.com/bshakr/ko/internal/validation"
	"github.com/spf13/cobra"
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
}

func runCleanup(cmd *cobra.Command, args []string) error {
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
	go func() {
		<-sigChan
		fmt.Println("\nOperation cancelled by user")
		cancel()
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

	// Close tmux window first (before removing worktree)
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
		}
	} else {
		fmt.Println("Not in a tmux session, skipping tmux cleanup")
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
