package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bshakr/ko/internal/config"
	"github.com/bshakr/ko/internal/git"
	"github.com/bshakr/ko/internal/tmux"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new <worktree-name>",
	Short: "Create a new worktree and tmux session",
	Long: `Create a new git worktree and automatically set up a tmux session with:
  - Top-left: vim
  - Bottom-left: setup script
  - Top-right: dev server (waits for setup)
  - Bottom-right: claude`,
	Args: cobra.ExactArgs(1),
	RunE: runNew,
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func runNew(cmd *cobra.Command, args []string) error {
	worktreeName := args[0]

	// Check if we're in a git repository
	if !git.IsGitRepo() {
		return fmt.Errorf("not in a git repository\nPlease run this command from within a git repository")
	}

	// Check if config exists, if not prompt user to run init
	exists, err := config.ConfigExists()
	if err != nil {
		return fmt.Errorf("failed to check for .koconfig: %w", err)
	}
	if !exists {
		return fmt.Errorf("no .koconfig found\nPlease run 'ko init' to set up your configuration first")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Determine the main repo root (handles both main repo and worktrees)
	var mainRepoRoot string
	if git.IsInWorktree() {
		mainRepoRoot, err = git.GetMainRepoRoot()
		if err != nil {
			return fmt.Errorf("failed to get main repository root: %w", err)
		}
	} else {
		mainRepoRoot, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	// Check if setup script exists (relative to main repo)
	setupPath := filepath.Join(mainRepoRoot, cfg.SetupScript)
	if _, err := os.Stat(setupPath); os.IsNotExist(err) {
		return fmt.Errorf("%s not found\nPlease create a setup script at %s", cfg.SetupScript, cfg.SetupScript)
	}

	// Check if dev script exists (relative to main repo)
	devPath := filepath.Join(mainRepoRoot, cfg.DevScript)
	if _, err := os.Stat(devPath); os.IsNotExist(err) {
		return fmt.Errorf("%s not found\nPlease create a dev script at %s", cfg.DevScript, cfg.DevScript)
	}

	// Create .ko directory if it doesn't exist
	koDir := filepath.Join(mainRepoRoot, ".ko")
	if err := os.MkdirAll(koDir, 0755); err != nil {
		return fmt.Errorf("failed to create .ko directory: %w", err)
	}

	// Check if worktree already exists
	worktreePath := filepath.Join(koDir, worktreeName)
	if _, err := os.Stat(worktreePath); err == nil {
		return fmt.Errorf("worktree .ko/%s already exists", worktreeName)
	}

	// Create git worktree
	fmt.Printf("Creating git worktree: .ko/%s\n", worktreeName)
	if err := git.CreateWorktree(worktreePath); err != nil {
		return fmt.Errorf("failed to create worktree: %w", err)
	}

	// Get repository name
	repoName, err := git.GetRepoName()
	if err != nil {
		return fmt.Errorf("failed to get repository name: %w", err)
	}

	// Create tmux session with config
	if err := tmux.CreateSession(repoName, worktreeName, worktreePath, cfg); err != nil {
		return fmt.Errorf("failed to create tmux session: %w", err)
	}

	fmt.Println("Worktree setup complete!")
	return nil
}
