package cmd

import (
	"fmt"
	"os"

	"github.com/bshakr/koh/internal/config"
	"github.com/bshakr/koh/internal/styles"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View current configuration",
	Long:  `Display the current koh configuration.`,
	RunE:  runConfig,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func runConfig(_ *cobra.Command, _ []string) error {
	// Get terminal width
	terminalWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || terminalWidth == 0 {
		terminalWidth = 80
	}

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
		Render(styles.RenderHelp("Run 'koh init' to change configuration interactively"))
	fmt.Println(help + "\n")

	return nil
}
