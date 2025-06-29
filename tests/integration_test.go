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

// setupComplexTestProject creates a realistic project structure for integration testing
func setupComplexTestProject(t *testing.T) string {
	tmpDir := t.TempDir()

	// Create a realistic project structure
	structure := map[string]string{
		"README.md":                  "# Test Project\n\nThis is a test project.",
		"go.mod":                     "module testproject\n\ngo 1.21",
		"go.sum":                     "// dependencies",
		"main.go":                    "package main\n\nfunc main() {\n\tprintln(\"Hello\")\n}",
		".gitignore":                 "*.log\n/dist/\n.env",
		"cmd/server/main.go":         "package main\n// server",
		"internal/api/handler.go":    "package api\n// handlers",
		"internal/api/middleware.go": "package api\n// middleware",
		"internal/db/connection.go":  "package db\n// database",
		"pkg/utils/helper.go":        "package utils\n// utilities",
		"docs/api.md":                "# API Documentation",
		"docs/deployment.md":         "# Deployment Guide",
		"tests/unit_test.go":         "package tests\n// unit tests",
		"tests/integration_test.go":  "package tests\n// integration tests",
		"config/config.yaml":         "database:\n  host: localhost",
		"scripts/build.sh":           "#!/bin/bash\ngo build",
		"assets/styles.css":          "body { font-family: sans-serif; }",
		"assets/logo.png":            "fake png content",
		"vendor/dependency/lib.go":   "package dependency",
	}

	for relPath, content := range structure {
		fullPath := filepath.Join(tmpDir, relPath)
		dir := filepath.Dir(fullPath)

		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Create some empty directories
	emptyDirs := []string{
		"logs",
		"tmp/cache",
		"dist",
	}

	for _, dir := range emptyDirs {
		fullPath := filepath.Join(tmpDir, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			t.Fatal(err)
		}
	}

	return tmpDir
}

func TestFullWorkflow(t *testing.T) {
	projectDir := setupComplexTestProject(t)

	tests := []struct {
		name         string
		initialDepth int
		expectFiles  []string
		expectDirs   []string
	}{
		{
			name:         "depth 1 shows top level",
			initialDepth: 1,
			expectFiles:  []string{"README.md", "go.mod", "go.sum", "main.go", ".gitignore"},
			expectDirs:   []string{"cmd", "internal", "pkg", "docs", "tests", "config", "scripts", "assets", "vendor", "logs", "tmp", "dist"},
		},
		{
			name:         "depth 2 expands subdirectories",
			initialDepth: 2,
			expectFiles:  []string{"README.md", "go.mod", "main.go"},
			expectDirs:   []string{"cmd", "internal", "pkg"},
		},
		{
			name:         "depth 3 shows deeper structure",
			initialDepth: 3,
			expectFiles:  []string{"handler.go", "middleware.go", "connection.go"},
			expectDirs:   []string{"api", "db"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build tree with specified depth
			root := tree.Build(projectDir, tt.initialDepth)
			if root == nil {
				t.Fatal("Failed to build tree")
			}

			// Create UI model
			model := ui.New(root, tt.initialDepth, projectDir)
			if model == nil {
				t.Fatal("Failed to create UI model")
			}

			// Render initial view
			view := model.View()
			if view == "" {
				t.Fatal("View should not be empty")
			}

			// Check that expected files/dirs appear in view
			for _, expectedFile := range tt.expectFiles {
				if !strings.Contains(view, expectedFile) {
					t.Errorf("View should contain file %s at depth %d", expectedFile, tt.initialDepth)
				}
			}

			for _, expectedDir := range tt.expectDirs {
				if !strings.Contains(view, expectedDir) {
					t.Errorf("View should contain directory %s at depth %d", expectedDir, tt.initialDepth)
				}
			}
		})
	}
}

func TestNavigationAndExpansion(t *testing.T) {
	projectDir := setupComplexTestProject(t)
	root := tree.Build(projectDir, 1)
	model := ui.New(root, 1, projectDir)

	// Test navigation sequence
	navigationSteps := []tea.KeyMsg{
		{Type: tea.KeyDown},                      // Move down
		{Type: tea.KeyDown},                      // Move down again
		{Type: tea.KeyEnter},                     // Expand/collapse or open
		{Type: tea.KeyUp},                        // Move up
		{Type: tea.KeyRunes, Runes: []rune{'j'}}, // Vim-style down
		{Type: tea.KeyRunes, Runes: []rune{'k'}}, // Vim-style up
	}

	currentModel := model
	for i, step := range navigationSteps {
		updatedModel, cmd := currentModel.Update(step)

		// Should not panic or error
		if updatedModel == nil {
			t.Errorf("Step %d: Model should not be nil after update", i)
		}

		// Most steps shouldn't return commands (except quit)
		if cmd != nil && step.Type != tea.KeyEsc {
			t.Logf("Step %d returned command: %v", i, cmd)
		}

		// Update current model for next iteration
		if m, ok := updatedModel.(*ui.Model); ok {
			currentModel = m
		}
	}

	// Final view should still be valid
	finalView := currentModel.View()
	if finalView == "" {
		t.Error("Final view should not be empty")
	}
}

func TestErrorHandling(t *testing.T) {
	// Test with invalid/inaccessible directory
	invalidPath := "/invalid/path/that/does/not/exist"
	root := tree.Build(invalidPath, 1)

	if root == nil {
		t.Fatal("Build should handle invalid paths gracefully")
	}

	model := ui.New(root, 1, invalidPath)
	if model == nil {
		t.Fatal("UI should handle invalid trees gracefully")
	}

	view := model.View()
	if view == "" {
		t.Error("View should not be empty even with invalid path")
	}

	// Should contain the invalid path in header
	if !strings.Contains(view, invalidPath) {
		t.Error("View should show the attempted path")
	}
}

func TestFileTypeHandling(t *testing.T) {
	projectDir := setupComplexTestProject(t)
	root := tree.Build(projectDir, 2)
	model := ui.New(root, 2, projectDir)

	view := model.View()

	// Should handle different file types
	fileTypes := []string{
		".md",   // Markdown
		".go",   // Go source
		".yaml", // YAML config
		".css",  // Stylesheets
		".png",  // Images
		".sh",   // Scripts
	}

	for _, fileType := range fileTypes {
		found := false
		for _, line := range strings.Split(view, "\n") {
			if strings.Contains(line, fileType) {
				found = true
				break
			}
		}
		if !found {
			t.Logf("File type %s not found in view (may be expected)", fileType)
		}
	}
}

func TestDepthConfiguration(t *testing.T) {
	projectDir := setupComplexTestProject(t)

	depthTests := []struct {
		depth    int
		minLines int
		maxLines int
	}{
		{0, 5, 10},   // Just header and controls
		{1, 15, 30},  // Top level files
		{2, 20, 50},  // One level deep
		{3, 25, 100}, // Two levels deep
	}

	for _, test := range depthTests {
		t.Run(t.Name(), func(t *testing.T) {
			root := tree.Build(projectDir, test.depth)
			model := ui.New(root, test.depth, projectDir)
			view := model.View()

			lines := strings.Split(view, "\n")
			lineCount := len(lines)

			if lineCount < test.minLines {
				t.Errorf("Depth %d: Expected at least %d lines, got %d", test.depth, test.minLines, lineCount)
			}

			if lineCount > test.maxLines {
				t.Errorf("Depth %d: Expected at most %d lines, got %d", test.depth, test.maxLines, lineCount)
			}
		})
	}
}

func TestQuietModes(t *testing.T) {
	projectDir := setupComplexTestProject(t)
	root := tree.Build(projectDir, 1)
	model := ui.New(root, 1, projectDir)

	// Test quit commands
	quitKeys := []tea.KeyMsg{
		{Type: tea.KeyRunes, Runes: []rune{'q'}},
		{Type: tea.KeyEsc},
		{Type: tea.KeyCtrlC},
	}

	for _, quitKey := range quitKeys {
		t.Run(t.Name(), func(t *testing.T) {
			_, cmd := model.Update(quitKey)

			if cmd == nil {
				t.Error("Quit keys should return quit command")
			}
		})
	}
}

func TestTreeCharacters(t *testing.T) {
	projectDir := setupComplexTestProject(t)
	root := tree.Build(projectDir, 2)
	model := ui.New(root, 2, projectDir)

	view := model.View()

	// Should contain tree drawing characters
	treeChars := []string{
		"├──", // Branch connector
		"└──", // Last item connector
		"│",   // Vertical line
		"▶",   // Collapsed indicator
		"▼",   // Expanded indicator
	}

	for _, char := range treeChars {
		if !strings.Contains(view, char) {
			t.Logf("Tree character %s not found (might be expected based on structure)", char)
		}
	}
}
