package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bshakr/ko/internal/git"
)

// Config represents the ko configuration
type Config struct {
	Editor      string   `json:"editor"`
	SetupScript string   `json:"setup_script"`
	DevScript   string   `json:"dev_script"`
	PaneCommands []string `json:"pane_commands"`
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		Editor:      "vim",
		SetupScript: "./bin/setup",
		DevScript:   "./bin/dev",
		PaneCommands: []string{
			"vim",
			"./bin/setup",
			"./bin/dev",
			"claude",
		},
	}
}

// ConfigPath returns the path to the .koconfig file in the repo root
func ConfigPath() (string, error) {
	// Check if we're in a git repository
	if !git.IsGitRepo() {
		return "", fmt.Errorf("not in a git repository")
	}

	// Get the main repo root (handles both main repo and worktrees)
	var repoRoot string
	var err error

	if git.IsInWorktree() {
		repoRoot, err = git.GetMainRepoRoot()
		if err != nil {
			return "", fmt.Errorf("failed to get main repository root: %w", err)
		}
	} else {
		repoRoot, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	return filepath.Join(repoRoot, ".koconfig"), nil
}

// ConfigExists checks if a .koconfig file exists in the repo
func ConfigExists() (bool, error) {
	configPath, err := ConfigPath()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// Load loads the configuration from disk
func Load() (*Config, error) {
	configPath, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	// If config doesn't exist, return an error
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("no .koconfig found - run 'ko init' to set up configuration")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// Save saves the configuration to disk
func (c *Config) Save() error {
	configPath, err := ConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Setup runs an interactive setup to create a .koconfig file
func Setup() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Ko Configuration Setup")
	fmt.Println("======================")
	fmt.Println()

	config := DefaultConfig()

	// Ask for editor
	fmt.Printf("Editor (default: %s): ", config.Editor)
	editor, _ := reader.ReadString('\n')
	editor = strings.TrimSpace(editor)
	if editor != "" {
		config.Editor = editor
	}

	// Ask for setup script
	fmt.Printf("Setup script (default: %s): ", config.SetupScript)
	setupScript, _ := reader.ReadString('\n')
	setupScript = strings.TrimSpace(setupScript)
	if setupScript != "" {
		config.SetupScript = setupScript
	}

	// Ask for dev script
	fmt.Printf("Dev script (default: %s): ", config.DevScript)
	devScript, _ := reader.ReadString('\n')
	devScript = strings.TrimSpace(devScript)
	if devScript != "" {
		config.DevScript = devScript
	}

	// Ask for pane commands
	fmt.Println()
	fmt.Println("Tmux pane commands (press Enter on empty line to finish):")
	fmt.Println("Current defaults:", config.PaneCommands)
	fmt.Print("Use defaults? (y/n): ")
	useDefaults, _ := reader.ReadString('\n')
	useDefaults = strings.TrimSpace(strings.ToLower(useDefaults))

	if useDefaults != "y" && useDefaults != "yes" {
		var paneCommands []string
		fmt.Println("Enter commands (one per line, empty line to finish):")
		for {
			fmt.Print("> ")
			cmd, _ := reader.ReadString('\n')
			cmd = strings.TrimSpace(cmd)
			if cmd == "" {
				break
			}
			paneCommands = append(paneCommands, cmd)
		}
		if len(paneCommands) > 0 {
			config.PaneCommands = paneCommands
		}
	}

	// Save the configuration
	if err := config.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	configPath, _ := ConfigPath()
	fmt.Println()
	fmt.Printf("Configuration saved to: %s\n", configPath)

	return nil
}
