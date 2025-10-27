package cmd

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestListModelInit(t *testing.T) {
	m := listModel{
		worktrees: []worktreeItem{
			{name: "test1", branch: "main", path: "/path/1", isCurrent: false},
		},
		inTmux: true,
	}

	cmd := m.Init()
	if cmd != nil {
		t.Error("Expected Init() to return nil")
	}
}

func TestListModelNavigationDown(t *testing.T) {
	worktrees := []worktreeItem{
		{name: "test1", branch: "main", path: "/path/1", isCurrent: false},
		{name: "test2", branch: "feature", path: "/path/2", isCurrent: true},
		{name: "test3", branch: "dev", path: "/path/3", isCurrent: false},
	}

	m := listModel{
		worktrees: worktrees,
		cursor:    0,
		inTmux:    true,
	}

	// Test: Down arrow navigation
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updatedModel.(listModel)
	if m.cursor != 1 {
		t.Errorf("Expected cursor at 1 after down arrow, got %d", m.cursor)
	}

	// Test: j (vim down) navigation
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updatedModel.(listModel)
	if m.cursor != 2 {
		t.Errorf("Expected cursor at 2 after j, got %d", m.cursor)
	}

	// Test: Can't go past the end
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updatedModel.(listModel)
	if m.cursor != 2 {
		t.Errorf("Expected cursor to stay at 2 (max), got %d", m.cursor)
	}
}

func TestListModelNavigationUp(t *testing.T) {
	worktrees := []worktreeItem{
		{name: "test1", branch: "main", path: "/path/1", isCurrent: false},
		{name: "test2", branch: "feature", path: "/path/2", isCurrent: true},
		{name: "test3", branch: "dev", path: "/path/3", isCurrent: false},
	}

	m := listModel{
		worktrees: worktrees,
		cursor:    2,
		inTmux:    true,
	}

	// Test: Up arrow navigation
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updatedModel.(listModel)
	if m.cursor != 1 {
		t.Errorf("Expected cursor at 1 after up arrow, got %d", m.cursor)
	}

	// Test: k (vim up) navigation
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updatedModel.(listModel)
	if m.cursor != 0 {
		t.Errorf("Expected cursor at 0 after k, got %d", m.cursor)
	}

	// Test: Can't go below 0
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updatedModel.(listModel)
	if m.cursor != 0 {
		t.Errorf("Expected cursor to stay at 0, got %d", m.cursor)
	}
}

func TestListModelJumpToTop(t *testing.T) {
	worktrees := []worktreeItem{
		{name: "test1", branch: "main", path: "/path/1", isCurrent: false},
		{name: "test2", branch: "feature", path: "/path/2", isCurrent: false},
		{name: "test3", branch: "dev", path: "/path/3", isCurrent: false},
	}

	m := listModel{
		worktrees: worktrees,
		cursor:    2,
		inTmux:    true,
	}

	// Test: g jumps to top
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
	m = updatedModel.(listModel)
	if m.cursor != 0 {
		t.Errorf("Expected cursor at 0 after g, got %d", m.cursor)
	}

	// Test: home also jumps to top
	m.cursor = 2
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyHome})
	m = updatedModel.(listModel)
	if m.cursor != 0 {
		t.Errorf("Expected cursor at 0 after home, got %d", m.cursor)
	}
}

func TestListModelJumpToBottom(t *testing.T) {
	worktrees := []worktreeItem{
		{name: "test1", branch: "main", path: "/path/1", isCurrent: false},
		{name: "test2", branch: "feature", path: "/path/2", isCurrent: false},
		{name: "test3", branch: "dev", path: "/path/3", isCurrent: false},
	}

	m := listModel{
		worktrees: worktrees,
		cursor:    0,
		inTmux:    true,
	}

	// Test: G jumps to bottom
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
	m = updatedModel.(listModel)
	if m.cursor != 2 {
		t.Errorf("Expected cursor at 2 after G, got %d", m.cursor)
	}

	// Test: end also jumps to bottom
	m.cursor = 0
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnd})
	m = updatedModel.(listModel)
	if m.cursor != 2 {
		t.Errorf("Expected cursor at 2 after end, got %d", m.cursor)
	}
}

func TestListModelSelection(t *testing.T) {
	worktrees := []worktreeItem{
		{name: "test1", branch: "main", path: "/path/1", isCurrent: false},
		{name: "test2", branch: "feature", path: "/path/2", isCurrent: false},
	}

	m := listModel{
		worktrees: worktrees,
		cursor:    1,
		inTmux:    true,
	}

	// Test: Enter selects worktree
	updatedModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updatedModel.(listModel)

	if m.selected != "test2" {
		t.Errorf("Expected selected='test2', got %q", m.selected)
	}
	if !m.switchSuccess {
		t.Error("Expected switchSuccess=true")
	}
	if cmd == nil {
		t.Error("Expected Quit command")
	}
}

func TestListModelSelectionNotInTmux(t *testing.T) {
	worktrees := []worktreeItem{
		{name: "test1", branch: "main", path: "/path/1", isCurrent: false},
	}

	m := listModel{
		worktrees: worktrees,
		cursor:    0,
		inTmux:    false, // Not in tmux
	}

	// Test: Enter should not select when not in tmux
	updatedModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updatedModel.(listModel)

	if m.selected != "" {
		t.Errorf("Expected no selection when not in tmux, got %q", m.selected)
	}
	if m.switchSuccess {
		t.Error("Expected switchSuccess=false when not in tmux")
	}
	if cmd != nil {
		t.Error("Expected no command when not in tmux")
	}
}

func TestListModelQuit(t *testing.T) {
	m := listModel{
		worktrees: []worktreeItem{{name: "test", branch: "main", path: "/", isCurrent: false}},
		inTmux:    true,
	}

	tests := []struct {
		name string
		key  tea.KeyMsg
	}{
		{"q key", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}},
		{"esc key", tea.KeyMsg{Type: tea.KeyEsc}},
		{"ctrl+c", tea.KeyMsg{Type: tea.KeyCtrlC}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedModel, cmd := m.Update(tt.key)
			result := updatedModel.(listModel)

			if !result.quitting {
				t.Errorf("Expected quitting=true for %s", tt.name)
			}
			if cmd == nil {
				t.Errorf("Expected Quit command for %s", tt.name)
			}
		})
	}
}

func TestListModelView(t *testing.T) {
	worktrees := []worktreeItem{
		{name: "test1", branch: "main", path: "/path/1", isCurrent: false},
		{name: "test2", branch: "feature", path: "/path/2", isCurrent: true},
	}

	m := listModel{
		worktrees: worktrees,
		cursor:    0,
		inTmux:    true,
	}

	view := m.View()

	// Check that the view contains expected elements
	if view == "" {
		t.Error("Expected non-empty view")
	}
	if !contains(view, "Ko Worktrees") {
		t.Error("Expected view to contain title")
	}
	if !contains(view, "test1") {
		t.Error("Expected view to contain worktree name")
	}
	if !contains(view, "main") {
		t.Error("Expected view to contain branch name")
	}
	if !contains(view, "[current]") {
		t.Error("Expected view to contain current label")
	}
}

func TestListModelViewQuitting(t *testing.T) {
	m := listModel{
		worktrees: []worktreeItem{{name: "test", branch: "main", path: "/", isCurrent: false}},
		quitting:  true,
		inTmux:    true,
	}

	view := m.View()

	// When quitting without switch success, view should be empty
	if view != "" {
		t.Errorf("Expected empty view when quitting, got %q", view)
	}
}

func TestListModelViewNotInTmux(t *testing.T) {
	m := listModel{
		worktrees: []worktreeItem{{name: "test", branch: "main", path: "/", isCurrent: false}},
		inTmux:    false,
	}

	view := m.View()

	// Should show different help text when not in tmux
	if !contains(view, "not in tmux") {
		t.Error("Expected view to indicate not in tmux")
	}
}

func TestListModelEmptyWorktreesList(t *testing.T) {
	m := listModel{
		worktrees: []worktreeItem{},
		cursor:    0,
		inTmux:    true,
	}

	// Navigation should handle empty list gracefully
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	result := updatedModel.(listModel)

	if result.cursor != 0 {
		t.Errorf("Expected cursor to stay at 0 with empty list, got %d", result.cursor)
	}

	// G should handle empty list
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
	result = updatedModel.(listModel)

	// With empty list, cursor should stay at 0 or be set to -1, but the check prevents issues
	if result.cursor < 0 {
		t.Errorf("Expected cursor >= 0 with empty list, got %d", result.cursor)
	}
}

// Helper function for substring checking
func contains(s, substr string) bool {
	return len(s) >= len(substr) && stringContains(s, substr)
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
