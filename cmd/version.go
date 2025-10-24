package cmd

import (
	"fmt"
	"os"

	"github.com/bshakr/ko/internal/styles"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// Version is the current version of ko
// This can be overridden at build time using ldflags:
// go build -ldflags="-X github.com/bshakr/ko/cmd.Version=v1.0.0"
var Version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display ko version",
	Long:  `Display the current version of ko.`,
	RunE:  runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) error {
	// Get terminal width
	terminalWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || terminalWidth == 0 {
		terminalWidth = 80
	}

	// Create version display
	versionText := fmt.Sprintf("ko version %s", Version)
	styledVersion := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Primary).
		Render(versionText)

	// Center the output
	centered := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(terminalWidth).
		Render(styledVersion)

	fmt.Println("\n" + centered + "\n")

	return nil
}
