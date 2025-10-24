// Package cmd implements the CLI commands for ko.
//
// Ko is a tool for managing git worktrees with automatic tmux session setup.
// It provides commands to create, list, and clean up worktrees with pre-configured
// development environments.
//
// The main commands are:
//   - new: Create a new worktree with a tmux session
//   - cleanup: Remove a worktree and close its tmux session
//   - list: Display all ko-managed worktrees
//   - init: Interactive configuration wizard
//   - config: Display current configuration
//
// Each command is implemented in its own file (new.go, cleanup.go, etc.).
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bshakr/ko/internal/config"
	"github.com/bshakr/ko/internal/git"
	"github.com/bshakr/ko/internal/styles"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var rootCmd = &cobra.Command{
	Use:   "ko",
	Short: "Git Worktree tmux Automation",
	Long: `ko - Git Worktree tmux Automation

A tool for managing git worktrees with automatic tmux session setup.
Creates isolated development environments with pre-configured panes.`,
	Run: runRoot,
}

func runRoot(_ *cobra.Command, _ []string) {
	// Get actual terminal width
	terminalWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || terminalWidth == 0 {
		terminalWidth = 80 // fallback to default
	}

	// Print large ASCII title
	asciiTitle := `
██╗  ██╗ ██████╗
██║ ██╔╝██╔═══██╗
█████╔╝ ██║   ██║
██╔═██╗ ██║   ██║
██║  ██╗╚██████╔╝
╚═╝  ╚═╝ ╚═════╝ `

	koTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Primary).
		Align(lipgloss.Center).
		Width(terminalWidth).
		Render(asciiTitle)

	subtitle := lipgloss.NewStyle().
		Foreground(styles.Subtle).
		Align(lipgloss.Center).
		Width(terminalWidth).
		Render("Git Worktree Manager")

	fmt.Println(koTitle)
	fmt.Println(subtitle)
	fmt.Println()

	// Check if in git repo
	if !git.IsGitRepo() {
		errorMsg := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(terminalWidth).
			Render(styles.ErrorMessage.Render("Not in a git repository"))

		helpMsg := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(terminalWidth).
			Render(styles.Muted.Render("Please run ko from within a git repository"))

		fmt.Println(errorMsg)
		fmt.Println(helpMsg)
		fmt.Println()
		return
	}

	// Get repository info
	repoName, _ := git.GetRepoName()

	// Get worktree count and current worktree
	worktreeCount := 0
	currentWorktree := ""

	var mainRepoRoot string
	var currentWorktreePath string

	if git.IsInWorktree() {
		mainRepoRoot, _ = git.GetMainRepoRoot()
		currentWorktreePath, _ = git.GetCurrentWorktreePath()
		currentWorktree = filepath.Base(currentWorktreePath)
	} else {
		mainRepoRoot, _ = os.Getwd()
		currentWorktree = "main"
	}

	// Count worktrees
	if mainRepoRoot != "" {
		koDir := filepath.Join(mainRepoRoot, ".ko")
		if _, err := os.Stat(koDir); err == nil {
			gitCmd := exec.Command("git", "worktree", "list")
			output, err := gitCmd.Output()
			if err == nil {
				lines := strings.Split(string(output), "\n")
				for _, line := range lines {
					if strings.Contains(line, "/.ko/") {
						worktreeCount++
					}
				}
			}
		}
	}

	// Check config status
	configExists, _ := config.ConfigExists()
	configStatus := styles.ErrorMessage.Render(styles.IconCross + " Not configured")
	if configExists {
		configStatus = styles.SuccessMessage.Render(styles.IconCheck + " Configured")
	}

	// Build status section - simple left-aligned, compact
	var statusContent strings.Builder
	statusContent.WriteString(styles.RenderKeyValue("Repository", repoName) + "\n")
	statusContent.WriteString(styles.RenderKeyValue("Worktrees", fmt.Sprintf("%d active", worktreeCount)) + "\n")
	statusContent.WriteString(styles.RenderKeyValue("Current", currentWorktree) + "\n")
	statusContent.WriteString(styles.Key.Render("Config:") + " " + configStatus)

	statusBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Subtle).
		Padding(0, 1).
		Render(statusContent.String())

	// Center the status box
	centeredStatusBox := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(terminalWidth).
		Render(statusBox)

	fmt.Println(centeredStatusBox)
	fmt.Println()

	// Quick Start section
	quickStartTitle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(terminalWidth).
		Render(styles.Subtitle.Render("Quick Start"))
	fmt.Println(quickStartTitle)

	quickStart := []struct {
		command string
		desc    string
	}{
		{"ko new <name>", "Create a new worktree"},
		{"ko list", "View all worktrees"},
		{"ko config", "Show configuration"},
	}

	// Find max command width for alignment
	maxCmdWidth := 0
	for _, qs := range quickStart {
		if len(qs.command) > maxCmdWidth {
			maxCmdWidth = len(qs.command)
		}
	}

	// Style for command column (colored but no background)
	cmdStyle := lipgloss.NewStyle().Foreground(styles.Warning)

	// Find max line length for this section
	maxLineLen := 0
	for _, qs := range quickStart {
		lineLen := maxCmdWidth + 3 + len(qs.desc) // cmd + spacing + desc
		if lineLen > maxLineLen {
			maxLineLen = lineLen
		}
	}

	for _, qs := range quickStart {
		// Manually pad the command to max width
		paddedCmd := fmt.Sprintf("%-*s", maxCmdWidth, qs.command)

		// Apply styling to the padded command
		styledCmd := cmdStyle.Render(paddedCmd)

		// Build the line with proper spacing
		line := styledCmd + "   " + qs.desc

		// Pad the entire line to max line length for consistent centering
		lineLenWithoutANSI := maxCmdWidth + 3 + len(qs.desc)
		paddingNeeded := maxLineLen - lineLenWithoutANSI
		if paddingNeeded > 0 {
			line = line + strings.Repeat(" ", paddingNeeded)
		}

		// Center the entire line
		centered := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(terminalWidth).
			Render(line)
		fmt.Println(centered)
	}
	fmt.Println()

	// Common Workflows section
	workflowsTitle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(terminalWidth).
		Render(styles.Subtitle.Render("Common Workflows"))
	fmt.Println(workflowsTitle)

	workflows := []struct {
		name    string
		command string
	}{
		{"Start new feature", "ko new feature-name"},
		{"List all worktrees", "ko list"},
		{"Clean up old work", "ko cleanup <name>"},
	}

	// Find max workflow name width for alignment
	maxNameWidth := 0
	for _, wf := range workflows {
		if len(wf.name) > maxNameWidth {
			maxNameWidth = len(wf.name)
		}
	}

	// Find max line length for this section
	maxWorkflowLineLen := 0
	for _, wf := range workflows {
		lineLen := maxNameWidth + 3 + len(wf.command)
		if lineLen > maxWorkflowLineLen {
			maxWorkflowLineLen = lineLen
		}
	}

	for _, wf := range workflows {
		// Manually pad the workflow name to max width
		paddedName := fmt.Sprintf("%-*s", maxNameWidth, wf.name)

		// Apply styling
		styledName := styles.Key.Render(paddedName)
		styledCommand := styles.Muted.Render(wf.command)

		// Build the line with proper spacing
		line := styledName + "   " + styledCommand

		// Pad the entire line to max line length for consistent centering
		lineLenWithoutANSI := maxNameWidth + 3 + len(wf.command)
		paddingNeeded := maxWorkflowLineLen - lineLenWithoutANSI
		if paddingNeeded > 0 {
			line = line + strings.Repeat(" ", paddingNeeded)
		}

		// Center the entire line
		centered := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(terminalWidth).
			Render(line)
		fmt.Println(centered)
	}
	fmt.Println()

	// Commands section
	commandsTitle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(terminalWidth).
		Render(styles.Subtitle.Render("Commands"))
	fmt.Println(commandsTitle)

	cmdGroups := []struct {
		title string
		cmds  []struct {
			name string
			desc string
		}
	}{
		{
			"Worktree Management",
			[]struct {
				name string
				desc string
			}{
				{"new", "Create new worktree + tmux session"},
				{"list", "List all worktrees"},
				{"cleanup", "Remove worktree and close session"},
			},
		},
		{
			"Configuration",
			[]struct {
				name string
				desc string
			}{
				{"init", "Interactive setup wizard"},
				{"config", "View current configuration"},
			},
		},
		{
			"Help",
			[]struct {
				name string
				desc string
			}{
				{"help", "Show help for any command"},
			},
		},
	}

	// Find max command name width across all groups for consistent alignment
	maxCmdNameWidth := 0
	for _, group := range cmdGroups {
		for _, cmd := range group.cmds {
			if len(cmd.name) > maxCmdNameWidth {
				maxCmdNameWidth = len(cmd.name)
			}
		}
	}

	// Find max line length using the padded command width
	maxCmdLineLen := 0
	for _, group := range cmdGroups {
		for _, cmd := range group.cmds {
			lineLen := maxCmdNameWidth + 3 + len(cmd.desc)
			if lineLen > maxCmdLineLen {
				maxCmdLineLen = lineLen
			}
		}
	}

	for i, group := range cmdGroups {
		if i > 0 {
			fmt.Println()
		}
		groupTitle := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(terminalWidth).
			Render(styles.Active.Render(group.title))
		fmt.Println(groupTitle)

		for _, cmd := range group.cmds {
			// Manually pad command name to max width
			paddedName := fmt.Sprintf("%-*s", maxCmdNameWidth, cmd.name)

			// Apply styling to command name (highlighted)
			styledCmdName := styles.Key.Render(paddedName)

			// Build the line with proper spacing
			line := styledCmdName + "   " + styles.Muted.Render(cmd.desc)

			// Pad the entire line to max line length for consistent centering
			lineLenWithoutANSI := maxCmdNameWidth + 3 + len(cmd.desc)
			paddingNeeded := maxCmdLineLen - lineLenWithoutANSI
			if paddingNeeded > 0 {
				line = line + strings.Repeat(" ", paddingNeeded)
			}

			// Center the line
			centered := lipgloss.NewStyle().
				Align(lipgloss.Center).
				Width(terminalWidth).
				Render(line)
			fmt.Println(centered)
		}
	}

	fmt.Println()

	// Context-aware tip
	var tip string
	if !configExists {
		tip = "💡 Tip: Run 'ko init' to set up your configuration first"
	} else if worktreeCount == 0 {
		tip = "💡 Tip: Run 'ko new feature-name' to create your first worktree"
	} else {
		tip = "💡 Tip: Use 'ko list' to see all your worktrees"
	}

	centeredTip := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(terminalWidth).
		Render(styles.Help.Render(tip))
	fmt.Println(centeredTip)
	fmt.Println()
}

// Execute runs the root command and handles any errors.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags can be added here if needed
}
