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

func runRoot(cmd *cobra.Command, args []string) {
	// Get actual terminal width
	terminalWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || terminalWidth == 0 {
		terminalWidth = 80 // fallback to default
	}

	// Use a max content width for better readability on wide screens
	// but make it responsive for narrow panes
	maxContentWidth := 70
	if terminalWidth < maxContentWidth {
		maxContentWidth = terminalWidth - 4 // leave some margin
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

	quickStart := []string{
		fmt.Sprintf("%s  Create a new worktree", styles.Code.Render("ko new <name>")),
		fmt.Sprintf("%s  View all worktrees", styles.Code.Render("ko list")),
		fmt.Sprintf("%s  Show configuration", styles.Code.Render("ko config")),
	}

	for _, item := range quickStart {
		centered := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(terminalWidth).
			Render(item)
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

	for _, wf := range workflows {
		line := fmt.Sprintf("%s  %s",
			styles.Key.Render(wf.name),
			styles.Muted.Render(wf.command))
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
		cmds  []string
	}{
		{
			"Worktree Management",
			[]string{
				"new        Create new worktree + tmux session",
				"list       List all worktrees",
				"cleanup    Remove worktree and close session",
			},
		},
		{
			"Configuration",
			[]string{
				"init       Interactive setup wizard",
				"config     View current configuration",
			},
		},
		{
			"Help",
			[]string{
				"help       Show help for any command",
			},
		},
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

		for _, cmdLine := range group.cmds {
			centered := lipgloss.NewStyle().
				Align(lipgloss.Center).
				Width(terminalWidth).
				Render(styles.Muted.Render(cmdLine))
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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags can be added here if needed
}
