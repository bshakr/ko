package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bshakr/ko/internal/git"
	"github.com/bshakr/ko/internal/styles"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all ko worktrees",
	Long:  `List all git worktrees in the .ko directory.`,
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(_ *cobra.Command, _ []string) error {
	// Get terminal width
	terminalWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || terminalWidth == 0 {
		terminalWidth = 80
	}

	// Check if we're in a git repository
	if !git.IsGitRepo() {
		return fmt.Errorf("not in a git repository")
	}

	// Determine the main repo root (handle being inside a worktree)
	var mainRepoRoot string
	var currentWorktreePath string

	if git.IsInWorktree() {
		mainRepoRoot, err = git.GetMainRepoRoot()
		if err != nil {
			return fmt.Errorf("failed to get main repository root: %w", err)
		}
		currentWorktreePath, err = git.GetCurrentWorktreePath()
		if err != nil {
			return fmt.Errorf("failed to get current worktree path: %w", err)
		}
	} else {
		mainRepoRoot, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	// Check if .ko directory exists
	koDir := filepath.Join(mainRepoRoot, ".ko")
	if _, err := os.Stat(koDir); os.IsNotExist(err) {
		msg := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(terminalWidth).
			Render(styles.Muted.Render("No worktrees found (no .ko directory)"))
		fmt.Println(msg)
		return nil
	}

	// List git worktrees
	gitCmd := exec.Command("git", "worktree", "list")
	output, err := gitCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to list worktrees: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	var koWorktrees []string

	for _, line := range lines {
		if strings.Contains(line, "/.ko/") {
			koWorktrees = append(koWorktrees, line)
		}
	}

	if len(koWorktrees) == 0 {
		msg := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(terminalWidth).
			Render(styles.Muted.Render("No ko worktrees found"))
		fmt.Println(msg)
		return nil
	}

	// Print title (centered)
	title := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(terminalWidth).
		Render(styles.RenderTitle(styles.IconTree + " Ko Worktrees"))
	fmt.Println("\n" + title)

	// Print worktrees (centered)
	for _, worktree := range koWorktrees {
		// Extract worktree name and branch
		parts := strings.Fields(worktree)
		if len(parts) >= 3 {
			path := parts[0]
			branch := strings.Trim(parts[len(parts)-1], "[]")
			name := filepath.Base(path)

			var line string
			// Highlight the current worktree
			if currentWorktreePath != "" && path == currentWorktreePath {
				icon := styles.Active.Render(styles.IconCurrent)
				nameStyled := styles.Active.Render(name)
				branchStyled := styles.HighlightStyle.Render(styles.IconBranch + " " + branch)
				currentLabel := styles.Muted.Render("[current]")
				line = fmt.Sprintf("%s %s %s %s", icon, nameStyled, branchStyled, currentLabel)
			} else {
				icon := styles.Muted.Render(styles.IconBullet)
				nameStyled := name
				branchStyled := styles.Muted.Render(styles.IconBranch + " " + branch)
				line = fmt.Sprintf("%s %s %s", icon, nameStyled, branchStyled)
			}

			centered := lipgloss.NewStyle().
				Align(lipgloss.Center).
				Width(terminalWidth).
				Render(line)
			fmt.Println(centered)
		}
	}

	// Print help text (centered)
	help := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(terminalWidth).
		Render(styles.RenderHelp("Use 'ko new <name>' to create a new worktree"))
	fmt.Println("\n" + help + "\n")

	return nil
}
