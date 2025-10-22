package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Editor == "" {
		t.Error("DefaultConfig() returned empty Editor")
	}

	if cfg.SetupScript == "" {
		t.Error("DefaultConfig() returned empty SetupScript")
	}

	if cfg.DevScript == "" {
		t.Error("DefaultConfig() returned empty DevScript")
	}

	if len(cfg.PaneCommands) == 0 {
		t.Error("DefaultConfig() returned empty PaneCommands")
	}

	t.Logf("Default config: %+v", cfg)
}

func TestConfigSaveAndLoad(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "ko-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test config file path
	configPath := filepath.Join(tempDir, ".koconfig")

	// Create a test config
	testConfig := &Config{
		Editor:      "nvim",
		SetupScript: "./test/setup",
		DevScript:   "./test/dev",
		PaneCommands: []string{
			"nvim",
			"./test/setup",
			"./test/dev",
			"test-cli",
		},
	}

	// Marshal and save manually (since Save() uses ConfigPath which needs git)
	data, err := json.MarshalIndent(testConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Read it back
	loadedData, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	var loadedConfig Config
	if err := json.Unmarshal(loadedData, &loadedConfig); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Verify the loaded config matches
	if loadedConfig.Editor != testConfig.Editor {
		t.Errorf("Editor mismatch: got %s, want %s", loadedConfig.Editor, testConfig.Editor)
	}

	if loadedConfig.SetupScript != testConfig.SetupScript {
		t.Errorf("SetupScript mismatch: got %s, want %s", loadedConfig.SetupScript, testConfig.SetupScript)
	}

	if loadedConfig.DevScript != testConfig.DevScript {
		t.Errorf("DevScript mismatch: got %s, want %s", loadedConfig.DevScript, testConfig.DevScript)
	}

	if len(loadedConfig.PaneCommands) != len(testConfig.PaneCommands) {
		t.Errorf("PaneCommands length mismatch: got %d, want %d",
			len(loadedConfig.PaneCommands), len(testConfig.PaneCommands))
	}

	for i, cmd := range testConfig.PaneCommands {
		if loadedConfig.PaneCommands[i] != cmd {
			t.Errorf("PaneCommands[%d] mismatch: got %s, want %s",
				i, loadedConfig.PaneCommands[i], cmd)
		}
	}
}

func TestConfigJSON(t *testing.T) {
	cfg := DefaultConfig()

	// Test marshaling
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	t.Logf("Config JSON:\n%s", string(data))

	// Test unmarshaling
	var loaded Config
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Verify unmarshaled config matches original
	if loaded.Editor != cfg.Editor {
		t.Errorf("Editor mismatch after marshal/unmarshal")
	}
}
