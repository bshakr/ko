package cmd

import (
	"fmt"

	"github.com/bshakr/ko/internal/config"
	"github.com/bshakr/ko/internal/styles"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View current configuration",
	Long:  `Display the current ko configuration.`,
	RunE:  runConfig,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func runConfig(_ *cobra.Command, _ []string) error {
	// Get terminal width
	terminalWidth := styles.GetTerminalWidth()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	configPath, err := config.ConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	// Print title (centered)
	title := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(terminalWidth).
		Render(styles.RenderTitle(styles.IconConfig + " Configuration"))

	location := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(terminalWidth).
		Render(styles.Muted.Render("Location: " + configPath))

	fmt.Println("\n" + title)
	fmt.Println(location)
	fmt.Println()

	// Create a styled box for config values
	var content string
	content += styles.RenderKeyValue("Setup Script", cfg.SetupScript) + "\n"
	content += "\n"
	if len(cfg.PaneCommands) > 0 {
		content += styles.Key.Render("Pane Commands:") + "\n"
		for i, cmdStr := range cfg.PaneCommands {
			content += fmt.Sprintf("  %d. %s\n", i+1, styles.Key.Render(cmdStr))
		}
	} else {
		content += styles.Muted.Render("No pane commands configured") + "\n"
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Subtle).
		Padding(0, 1).
		Render(content)

	centeredBox := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(terminalWidth).
		Render(box)

	fmt.Println(centeredBox)
	fmt.Println()

	help := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(terminalWidth).
		Render(styles.RenderHelp("Run 'ko init' to change configuration interactively"))
	fmt.Println(help + "\n")

	return nil
}
