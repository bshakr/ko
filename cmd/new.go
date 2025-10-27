package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bshakr/ko/internal/config"
	"github.com/bshakr/ko/internal/git"
	"github.com/bshakr/ko/internal/signals"
	"github.com/bshakr/ko/internal/tmux"
	"github.com/bshakr/ko/internal/validation"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new <worktree-name>",
	Short: "Create a new worktree and tmux session",
	Long: `Create a new git worktree and automatically set up a tmux session.
The session will have one pane for the setup script and additional panes for configured commands.`,
	Args: cobra.ExactArgs(1),
	RunE: runNew,
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func runNew(_ *cobra.Command, args []string) error {
	worktreeName := args[0]

	// Validate worktree name for security
	if err := validation.ValidateWorktreeName(worktreeName); err != nil {
		return fmt.Errorf("invalid worktree name: %w", err)
	}

	// Set up context with cancellation for long-running operations and signal handling
	ctx, cleanup := signals.SetupCancellableContext()
	defer cleanup()

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
	mainRepoRoot, err := git.GetMainRepoRootOrCwd()
	if err != nil {
		return fmt.Errorf("failed to get repository root: %w", err)
	}

	// Check if setup script exists and is within repository boundaries
	if cfg.SetupScript != "" {
		setupPath := filepath.Join(mainRepoRoot, cfg.SetupScript)

		// Validate setup script is within repository (security check)
		if err := validation.ValidatePathWithinRepository(setupPath, mainRepoRoot); err != nil {
			return fmt.Errorf("setup script %w\nAttempted path: %s", err, cfg.SetupScript)
		}

		// Check if the script exists
		if _, err := os.Stat(setupPath); os.IsNotExist(err) {
			return fmt.Errorf("%s not found\nPlease create a setup script at %s", cfg.SetupScript, cfg.SetupScript)
		}
	}

	// Create .ko directory if it doesn't exist
	koDir := filepath.Join(mainRepoRoot, ".ko")
	//nolint:gosec // G301: 0755 is standard permission for user directories
	if err := os.MkdirAll(koDir, 0755); err != nil {
		return fmt.Errorf("failed to create .ko directory: %w", err)
	}

	// Check if worktree already exists
	worktreePath := filepath.Join(koDir, worktreeName)
	if _, err := os.Stat(worktreePath); err == nil {
		return fmt.Errorf("worktree .ko/%s already exists", worktreeName)
	}

	// Create git worktree with context
	fmt.Printf("Creating git worktree: .ko/%s\n", worktreeName)
	if err := git.CreateWorktreeWithContext(ctx, worktreePath); err != nil {
		return fmt.Errorf("failed to create worktree: %w", err)
	}

	// Get repository name
	repoName, err := git.GetRepoName()
	if err != nil {
		return fmt.Errorf("failed to get repository name: %w", err)
	}

	// Create tmux session with config and context
	if err := tmux.CreateSessionWithContext(ctx, repoName, worktreeName, worktreePath, cfg); err != nil {
		return fmt.Errorf("failed to create tmux session: %w", err)
	}

	fmt.Println("Worktree setup complete!")
	return nil
}
