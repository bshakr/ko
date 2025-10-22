package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bshakr/ko/internal/git"
	"github.com/bshakr/ko/internal/tmux"
	"github.com/spf13/cobra"
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup <worktree-name>",
	Short: "Close tmux session and remove worktree",
	Long:  `Closes the associated tmux window and removes the git worktree.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runCleanup,
}

func init() {
	rootCmd.AddCommand(cleanupCmd)
}

func runCleanup(cmd *cobra.Command, args []string) error {
	worktreeName := args[0]

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

	// Remove the git worktree
	if worktreeExists {
		fmt.Printf("Removing git worktree: .ko/%s\n", worktreeName)
		if err := git.RemoveWorktree(worktreePath); err != nil {
			fmt.Printf("Warning: Failed to remove worktree automatically: %v\n", err)
			fmt.Printf("You may need to run: git worktree remove .ko/%s --force\n", worktreeName)
			fmt.Println("Or manually delete uncommitted changes first")
		} else {
			fmt.Println("Worktree removed successfully")
		}
	}

	// Close tmux window
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

	fmt.Println("Cleanup complete!")
	return nil
}
