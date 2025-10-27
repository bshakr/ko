package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.SetupScript == "" {
		t.Error("DefaultConfig() returned empty SetupScript")
	}

	if cfg.PaneCommands == nil {
		t.Error("DefaultConfig() returned nil PaneCommands")
	}

	t.Logf("Default config: %+v", cfg)
}

func TestConfigSaveAndLoad(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "ko-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Create a test config file path
	configPath := filepath.Join(tempDir, ".kohconfig")

	// Create a test config
	testConfig := &Config{
		SetupScript: "./test/setup",
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

	//nolint:gosec // G306: Test file - 0644 is acceptable for temp test files
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Read it back
	//nolint:gosec // G304: Test file - reading test config is expected
	loadedData, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	var loadedConfig Config
	if err := json.Unmarshal(loadedData, &loadedConfig); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Verify the loaded config matches
	if loadedConfig.SetupScript != testConfig.SetupScript {
		t.Errorf("SetupScript mismatch: got %s, want %s", loadedConfig.SetupScript, testConfig.SetupScript)
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
	if loaded.SetupScript != cfg.SetupScript {
		t.Errorf("SetupScript mismatch after marshal/unmarshal")
	}
}
