// Package styles defines consistent styling for ko's terminal UI.
//
// This package provides:
//   - Color palette (Primary, Success, Warning, Error, etc.)
//   - Text styles (Title, Subtitle, Key, Value, Muted, etc.)
//   - Icons for visual consistency
//   - Helper functions for formatted output
//
// All terminal output should use these styles to maintain a consistent
// look and feel across the application. The styles use the lipgloss library
// for terminal-aware color rendering.
package styles

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

// Base colors using terminal theme colors (ANSI)
// These will adapt to the user's terminal color scheme
var (
	// Primary color - typically blue/cyan in most themes
	Primary = lipgloss.AdaptiveColor{Light: "4", Dark: "12"}

	// Success color - typically green
	Success = lipgloss.AdaptiveColor{Light: "2", Dark: "10"}

	// Warning color - typically yellow
	Warning = lipgloss.AdaptiveColor{Light: "3", Dark: "11"}

	// Error color - typically red
	Error = lipgloss.AdaptiveColor{Light: "1", Dark: "9"}

	// Muted/subtle color - typically gray
	Subtle = lipgloss.AdaptiveColor{Light: "8", Dark: "7"}

	// Highlight - typically magenta/purple
	Highlight = lipgloss.AdaptiveColor{Light: "5", Dark: "13"}
)

// Common styles
var (
	// Title style for headers
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(Primary).
		MarginBottom(1)

	// Subtitle for secondary headers
	Subtitle = lipgloss.NewStyle().
			Foreground(Primary).
			MarginBottom(1)

	// Active/Current item indicator
	Active = lipgloss.NewStyle().
		Bold(true).
		Foreground(Highlight)

	// Highlighted text
	HighlightStyle = lipgloss.NewStyle().
			Foreground(Highlight)

	// Success message
	SuccessMessage = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true)

	// Error message
	ErrorMessage = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true)

	// Subtle/muted text
	Muted = lipgloss.NewStyle().
		Foreground(Subtle)

	// Key-value pair key
	Key = lipgloss.NewStyle().
		Foreground(Primary).
		Bold(true)

	// Box/border style
	Box = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Subtle).
		Padding(1, 2)

	// List item style
	ListItem = lipgloss.NewStyle().
			PaddingLeft(2)

	// Active list item
	ActiveListItem = lipgloss.NewStyle().
			PaddingLeft(0).
			Bold(true).
			Foreground(Highlight)

	// Code/inline style
	Code = lipgloss.NewStyle().
		Foreground(Warning).
		Background(lipgloss.AdaptiveColor{Light: "0", Dark: "8"}).
		Padding(0, 1)

	// Help text
	Help = lipgloss.NewStyle().
		Foreground(Subtle).
		Italic(true)
)

// Symbols using Unicode that work well in most terminals
const (
	IconCheck   = "✓"
	IconCross   = "✗"
	IconArrow   = "→"
	IconBullet  = "•"
	IconCurrent = "❯"
	IconConfig  = "⚙"
	IconBranch  = "⎇"
	IconTree    = "⚘"
)

// RenderTitle renders text with the Title style.
func RenderTitle(text string) string {
	return Title.Render(text)
}

// RenderSubtitle renders text with the Subtitle style.
func RenderSubtitle(text string) string {
	return Subtitle.Render(text)
}

// RenderSuccess renders text as a success message with a check icon.
func RenderSuccess(text string) string {
	return SuccessMessage.Render(IconCheck + " " + text)
}

// RenderError renders text as an error message with a cross icon.
func RenderError(text string) string {
	return ErrorMessage.Render(IconCross + " " + text)
}

// RenderKeyValue renders a key-value pair with styled key.
func RenderKeyValue(key, value string) string {
	return Key.Render(key+":") + " " + value
}

// RenderBox renders content inside a styled box.
func RenderBox(content string) string {
	return Box.Render(content)
}

// RenderHelp renders text with the Help style.
func RenderHelp(text string) string {
	return Help.Render(text)
}

// GetTerminalWidth returns the current terminal width, defaulting to 80 if unavailable.
// This is useful for centering and formatting output to fit the terminal.
func GetTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width == 0 {
		return 80
	}
	return width
}
