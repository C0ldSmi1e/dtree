package tests

import (
	"dtree/internal/tree"
	"dtree/internal/ui"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func createTestTree(t *testing.T) (*tree.Node, string) {
	tmpDir := t.TempDir()

	// Create simple test structure
	files := map[string]string{
		"file1.txt":         "content1",
		"file2.go":          "package main",
		"subdir/nested.txt": "nested",
	}

	for relPath, content := range files {
		fullPath := filepath.Join(tmpDir, relPath)
		dir := filepath.Dir(fullPath)

		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	root := tree.Build(tmpDir, 2) // Build with depth 2
	return root, tmpDir
}

func TestUIModelCreation(t *testing.T) {
	root, rootPath := createTestTree(t)

	model := ui.New(root, 2, rootPath)

	if model == nil {
		t.Fatal("ui.New returned nil")
	}

	// Test basic model properties
	if model == nil {
		t.Error("Model should not be nil")
	}
}

func TestUIModelInit(t *testing.T) {
	root, rootPath := createTestTree(t)
	model := ui.New(root, 1, rootPath)

	cmd := model.Init()
	if cmd != nil {
		t.Error("Init should return nil cmd")
	}
}

func TestUIModelView(t *testing.T) {
	root, rootPath := createTestTree(t)
	model := ui.New(root, 1, rootPath)

	view := model.View()

	// Basic view validation
	if view == "" {
		t.Error("View should not be empty")
	}

	// Should contain header
	if !strings.Contains(view, "DTree") {
		t.Error("View should contain DTree header")
	}

	// Should contain controls
	if !strings.Contains(view, "Controls:") {
		t.Error("View should contain controls information")
	}

	// Should contain root directory name
	rootName := filepath.Base(rootPath)
	if !strings.Contains(view, rootName) {
		t.Error("View should contain root directory name")
	}
}

func TestUIModelNavigation(t *testing.T) {
	root, rootPath := createTestTree(t)
	model := ui.New(root, 1, rootPath)

	// Test down navigation
	// initialCursor := 0

	// Simulate down key press
	downKey := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, cmd := model.Update(downKey)

	if cmd != nil {
		t.Error("Navigation should not return command")
	}

	// Can't easily check cursor position without accessing internals
	// This tests that Update doesn't panic
	if updatedModel == nil {
		t.Error("Update should return model")
	}
}

func TestUIModelKeyHandling(t *testing.T) {
	root, rootPath := createTestTree(t)
	model := ui.New(root, 1, rootPath)

	tests := []struct {
		name     string
		key      tea.KeyMsg
		wantQuit bool
	}{
		{
			name:     "q key quits",
			key:      tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
			wantQuit: true,
		},
		{
			name:     "esc key quits",
			key:      tea.KeyMsg{Type: tea.KeyEsc},
			wantQuit: true,
		},
		{
			name:     "ctrl+c quits",
			key:      tea.KeyMsg{Type: tea.KeyCtrlC},
			wantQuit: true,
		},
		{
			name:     "up key navigates",
			key:      tea.KeyMsg{Type: tea.KeyUp},
			wantQuit: false,
		},
		{
			name:     "down key navigates",
			key:      tea.KeyMsg{Type: tea.KeyDown},
			wantQuit: false,
		},
		{
			name:     "j key navigates",
			key:      tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
			wantQuit: false,
		},
		{
			name:     "k key navigates",
			key:      tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}},
			wantQuit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, cmd := model.Update(tt.key)

			if tt.wantQuit {
				if cmd == nil {
					t.Error("Quit keys should return quit command")
				}
				// Check if it's a quit command (this is a bit hacky)
				if cmd != nil && cmd() == nil {
					t.Error("Expected quit command")
				}
			} else {
				if cmd != nil {
					// Navigation keys shouldn't return commands
					t.Error("Navigation keys should not return commands")
				}
			}
		})
	}
}

func TestUIModelStatusMessage(t *testing.T) {
	root, rootPath := createTestTree(t)
	model := ui.New(root, 1, rootPath)

	// Test setting status
	model.SetStatus("Test error message")

	view := model.View()

	if !strings.Contains(view, "Test error message") {
		t.Error("View should contain status message")
	}

	// Test clearing status
	model.SetStatus("")

	view = model.View()

	if strings.Contains(view, "Test error message") {
		t.Error("View should not contain cleared status message")
	}
}

func TestUIModelExpandCollapse(t *testing.T) {
	root, rootPath := createTestTree(t)
	model := ui.New(root, 1, rootPath)

	// Find a directory in the view
	// initialView := model.View()

	// Simulate enter key to expand/collapse
	enterKey := tea.KeyMsg{Type: tea.KeyEnter}
	_, cmd := model.Update(enterKey)

	if cmd != nil {
		t.Error("Enter key should not return command for directory operations")
	}

	// Get view after enter
	afterView := model.View()

	// View should potentially change (expand/collapse)
	// This is a basic test that it doesn't panic
	if afterView == "" {
		t.Error("View should not be empty after enter")
	}
}

func TestUIModelFileOpening(t *testing.T) {
	root, rootPath := createTestTree(t)
	model := ui.New(root, 1, rootPath)

	// This test ensures file opening doesn't panic
	// We can't easily test actual file opening without system dependencies

	// Simulate selecting a file and pressing enter
	// (This is simplified since we can't easily navigate to a specific file)
	enterKey := tea.KeyMsg{Type: tea.KeyEnter}
	_, cmd := model.Update(enterKey)

	// Should not panic and should not return command
	if cmd != nil {
		t.Error("File opening should not return command")
	}
}

func TestUIModelNonKeyMessage(t *testing.T) {
	root, rootPath := createTestTree(t)
	model := ui.New(root, 1, rootPath)

	// Test with non-key message
	type customMsg struct{}

	updatedModel, cmd := model.Update(customMsg{})

	if cmd != nil {
		t.Error("Non-key messages should not return commands")
	}

	if updatedModel == nil {
		t.Error("Should return model even for unknown messages")
	}
}
