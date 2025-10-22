package cmd

import (
	"fmt"
	"strings"

	"github.com/bshakr/ko/internal/config"
	"github.com/bshakr/ko/internal/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Interactive configuration setup",
	Long:  `Run an interactive wizard to configure ko settings.`,
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

type step int

const (
	stepEditor step = iota
	stepSetupScript
	stepDevScript
	stepPane0
	stepPane1
	stepPane2
	stepPane3
	stepConfirm
	stepDone
)

type initModel struct {
	step     step
	config   *config.Config
	inputs   []textinput.Model
	focusIdx int
	err      error
}

func initialModel() initModel {
	cfg := config.DefaultConfig()

	inputs := make([]textinput.Model, 7)

	// Editor input
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "vim"
	inputs[0].SetValue(cfg.Editor)
	inputs[0].Focus()
	inputs[0].CharLimit = 50
	inputs[0].Width = 50
	inputs[0].Prompt = "❯ "

	// Setup script input
	inputs[1] = textinput.New()
	inputs[1].Placeholder = "./bin/setup"
	inputs[1].SetValue(cfg.SetupScript)
	inputs[1].CharLimit = 100
	inputs[1].Width = 50
	inputs[1].Prompt = "❯ "

	// Dev script input
	inputs[2] = textinput.New()
	inputs[2].Placeholder = "./bin/dev"
	inputs[2].SetValue(cfg.DevScript)
	inputs[2].CharLimit = 100
	inputs[2].Width = 50
	inputs[2].Prompt = "❯ "

	// Pane commands
	for i := 0; i < 4; i++ {
		inputs[3+i] = textinput.New()
		inputs[3+i].SetValue(cfg.PaneCommands[i])
		inputs[3+i].CharLimit = 100
		inputs[3+i].Width = 50
		inputs[3+i].Prompt = "❯ "
	}

	return initModel{
		step:     stepEditor,
		config:   cfg,
		inputs:   inputs,
		focusIdx: 0,
	}
}

func (m initModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m initModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			// Move to next step
			m.saveCurrentInput()

			if m.step == stepDone {
				return m, tea.Quit
			}

			m.step++
			if m.step <= stepPane3 {
				m.focusIdx = int(m.step)
				m.inputs[m.focusIdx].Focus()
			}
			return m, nil

		case "tab", "shift+tab", "up", "down":
			// Navigate between inputs in confirm step
			if m.step == stepConfirm {
				s := msg.String()
				if s == "up" || s == "shift+tab" {
					m.focusIdx--
				} else {
					m.focusIdx++
				}

				if m.focusIdx > 6 {
					m.focusIdx = 0
				} else if m.focusIdx < 0 {
					m.focusIdx = 6
				}

				for i := range m.inputs {
					if i == m.focusIdx {
						m.inputs[i].Focus()
					} else {
						m.inputs[i].Blur()
					}
				}
			}
			return m, nil
		}
	}

	// Update the focused input
	if m.step <= stepPane3 {
		var cmd tea.Cmd
		m.inputs[m.focusIdx], cmd = m.inputs[m.focusIdx].Update(msg)
		return m, cmd
	}

	// Update all inputs in confirm step
	if m.step == stepConfirm {
		var cmd tea.Cmd
		m.inputs[m.focusIdx], cmd = m.inputs[m.focusIdx].Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *initModel) saveCurrentInput() {
	switch m.step {
	case stepEditor:
		m.config.Editor = m.inputs[0].Value()
	case stepSetupScript:
		m.config.SetupScript = m.inputs[1].Value()
	case stepDevScript:
		m.config.DevScript = m.inputs[2].Value()
	case stepPane0:
		m.config.PaneCommands[0] = m.inputs[3].Value()
	case stepPane1:
		m.config.PaneCommands[1] = m.inputs[4].Value()
	case stepPane2:
		m.config.PaneCommands[2] = m.inputs[5].Value()
	case stepPane3:
		m.config.PaneCommands[3] = m.inputs[6].Value()
	case stepConfirm:
		// Save all values
		m.config.Editor = m.inputs[0].Value()
		m.config.SetupScript = m.inputs[1].Value()
		m.config.DevScript = m.inputs[2].Value()
		for i := 0; i < 4; i++ {
			m.config.PaneCommands[i] = m.inputs[3+i].Value()
		}
		// Save to disk
		if err := m.config.Save(); err != nil {
			m.err = err
		} else {
			m.step = stepDone
		}
	}
}

func (m initModel) View() string {
	var b strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Primary).
		MarginBottom(1).
		Render(styles.IconConfig + " Ko Configuration Setup")
	b.WriteString("\n" + title + "\n\n")

	switch m.step {
	case stepEditor:
		b.WriteString(styles.Subtitle.Render("Which editor would you like to use?"))
		b.WriteString("\n\n  ")
		b.WriteString(m.inputs[0].View())
		b.WriteString("\n\n")
		b.WriteString(styles.Muted.Render("  Examples: vim, nvim, code, emacs"))
		b.WriteString("\n")

	case stepSetupScript:
		b.WriteString(styles.Subtitle.Render("Path to your setup script?"))
		b.WriteString("\n\n  ")
		b.WriteString(m.inputs[1].View())
		b.WriteString("\n\n")
		b.WriteString(styles.Muted.Render("  This script runs once when creating a new worktree"))
		b.WriteString("\n")

	case stepDevScript:
		b.WriteString(styles.Subtitle.Render("Path to your dev script?"))
		b.WriteString("\n\n  ")
		b.WriteString(m.inputs[2].View())
		b.WriteString("\n\n")
		b.WriteString(styles.Muted.Render("  This script starts your development server"))
		b.WriteString("\n")

	case stepPane0:
		b.WriteString(styles.Subtitle.Render("Command for top-left pane?"))
		b.WriteString("\n\n  ")
		b.WriteString(m.inputs[3].View())
		b.WriteString("\n\n")
		b.WriteString(styles.Muted.Render("  Default: vim (or your chosen editor)"))
		b.WriteString("\n")

	case stepPane1:
		b.WriteString(styles.Subtitle.Render("Command for bottom-left pane?"))
		b.WriteString("\n\n  ")
		b.WriteString(m.inputs[4].View())
		b.WriteString("\n\n")
		b.WriteString(styles.Muted.Render("  Default: ./bin/setup"))
		b.WriteString("\n")

	case stepPane2:
		b.WriteString(styles.Subtitle.Render("Command for top-right pane?"))
		b.WriteString("\n\n  ")
		b.WriteString(m.inputs[5].View())
		b.WriteString("\n\n")
		b.WriteString(styles.Muted.Render("  Default: ./bin/dev (runs after setup completes)"))
		b.WriteString("\n")

	case stepPane3:
		b.WriteString(styles.Subtitle.Render("Command for bottom-right pane?"))
		b.WriteString("\n\n  ")
		b.WriteString(m.inputs[6].View())
		b.WriteString("\n\n")
		b.WriteString(styles.Muted.Render("  Default: claude"))
		b.WriteString("\n")

	case stepConfirm:
		b.WriteString(styles.Subtitle.Render("Review your configuration:"))
		b.WriteString("\n\n")

		var content string
		content += styles.RenderKeyValue("Editor", m.inputs[0].Value()) + "\n"
		content += styles.RenderKeyValue("Setup Script", m.inputs[1].Value()) + "\n"
		content += styles.RenderKeyValue("Dev Script", m.inputs[2].Value()) + "\n"
		content += "\n"
		content += styles.Key.Render("Pane Commands:") + "\n"
		content += fmt.Sprintf("  1. %s\n", styles.Code.Render(m.inputs[3].Value()))
		content += fmt.Sprintf("  2. %s\n", styles.Code.Render(m.inputs[4].Value()))
		content += fmt.Sprintf("  3. %s\n", styles.Code.Render(m.inputs[5].Value()))
		content += fmt.Sprintf("  4. %s\n", styles.Code.Render(m.inputs[6].Value()))

		box := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Subtle).
			Padding(1, 2).
			Render(content)

		b.WriteString(box)
		b.WriteString("\n\n")
		b.WriteString(styles.Muted.Render("  Use Tab/Shift+Tab to edit, Enter to save"))
		b.WriteString("\n")

	case stepDone:
		if m.err != nil {
			b.WriteString(styles.RenderError("Error saving configuration: " + m.err.Error()))
		} else {
			configPath, _ := config.ConfigPath()
			b.WriteString(styles.RenderSuccess("Configuration saved!"))
			b.WriteString("\n\n")
			b.WriteString(styles.Muted.Render("  Config location: " + configPath))
		}
		b.WriteString("\n")
	}

	if m.step < stepConfirm {
		b.WriteString("\n")
		b.WriteString(styles.Help.Render("  Press Enter to continue, Ctrl+C to cancel"))
		b.WriteString("\n")
	} else if m.step == stepConfirm {
		b.WriteString("\n")
		b.WriteString(styles.Help.Render("  Press Enter to save, Ctrl+C to cancel"))
		b.WriteString("\n")
	}

	return b.String()
}

func runInit(cmd *cobra.Command, args []string) error {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running interactive setup: %w", err)
	}
	return nil
}
