package tests

import (
	"dtree/internal/tree"
	"os"
	"path/filepath"
	"testing"
)

// setupTestFixture creates a test directory structure
func setupTestFixture(t *testing.T) string {
	tmpDir := t.TempDir()

	// Create test structure:
	// testdir/
	// ├── file1.txt
	// ├── file2.go
	// ├── .hidden
	// └── subdir/
	//     ├── nested.txt
	//     └── empty_dir/

	files := map[string]string{
		"file1.txt":         "sample content",
		"file2.go":          "package main\n\nfunc main() {}",
		".hidden":           "hidden file",
		"subdir/nested.txt": "nested content",
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

	// Create empty directory
	emptyDir := filepath.Join(tmpDir, "subdir", "empty_dir")
	if err := os.Mkdir(emptyDir, 0755); err != nil {
		t.Fatal(err)
	}

	return tmpDir
}

func TestTreeBuild(t *testing.T) {
	testDir := setupTestFixture(t)

	tests := []struct {
		name         string
		initialDepth int
		wantChildren int
		checkSubdir  bool
	}{
		{
			name:         "depth 0 loads no children",
			initialDepth: 0,
			wantChildren: 0,
			checkSubdir:  false,
		},
		{
			name:         "depth 1 loads immediate children",
			initialDepth: 1,
			wantChildren: 4, // file1.txt, file2.go, .hidden, subdir
			checkSubdir:  false,
		},
		{
			name:         "depth 2 expands subdirectories",
			initialDepth: 2,
			wantChildren: 4,
			checkSubdir:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := tree.Build(testDir, tt.initialDepth)

			// Basic validation
			if root == nil {
				t.Fatal("Build returned nil")
			}

			if !root.IsDir {
				t.Error("Root should be a directory")
			}

			if root.Depth != 0 {
				t.Errorf("Root depth = %d, want 0", root.Depth)
			}

			if root.Parent != nil {
				t.Error("Root should have no parent")
			}

			// Check children count
			if tt.initialDepth > 0 {
				if len(root.Children) != tt.wantChildren {
					t.Errorf("Children count = %d, want %d", len(root.Children), tt.wantChildren)
				}
			} else {
				if len(root.Children) != 0 {
					t.Errorf("Depth 0 should have no children, got %d", len(root.Children))
				}
			}

			// Check subdirectory expansion for depth 2
			if tt.checkSubdir {
				var subdir *tree.Node
				for _, child := range root.Children {
					if child.Name == "subdir" && child.IsDir {
						subdir = child
						break
					}
				}

				if subdir == nil {
					t.Fatal("Should have found subdir")
				}

				if !subdir.IsExpanded {
					t.Error("Subdir should be expanded at depth 2")
				}

				if len(subdir.Children) != 2 { // nested.txt, empty_dir
					t.Errorf("Subdir should have 2 children, got %d", len(subdir.Children))
				}
			}
		})
	}
}

func TestNodeLoadChildren(t *testing.T) {
	testDir := setupTestFixture(t)

	root := &tree.Node{
		Name:  filepath.Base(testDir),
		Path:  testDir,
		IsDir: true,
		Depth: 0,
	}

	// Initially no children
	if len(root.Children) != 0 {
		t.Error("Node should start with no children")
	}

	// Load children
	root.LoadChildren()

	// Should have 4 children
	if len(root.Children) != 4 {
		t.Errorf("LoadChildren should create 4 children, got %d", len(root.Children))
	}

	// Validate parent-child relationships
	for _, child := range root.Children {
		if child.Parent != root {
			t.Error("Child parent should point to root")
		}
		if child.Depth != 1 {
			t.Errorf("Child depth should be 1, got %d", child.Depth)
		}
	}

	// Check file vs directory detection
	fileCount := 0
	dirCount := 0
	for _, child := range root.Children {
		if child.IsDir {
			dirCount++
		} else {
			fileCount++
		}
	}

	if fileCount != 3 { // file1.txt, file2.go, .hidden
		t.Errorf("Expected 3 files, got %d", fileCount)
	}

	if dirCount != 1 { // subdir
		t.Errorf("Expected 1 directory, got %d", dirCount)
	}
}

func TestNodeIsLastChild(t *testing.T) {
	root := &tree.Node{Name: "root", IsDir: true, Depth: 0}
	child1 := &tree.Node{Name: "first", Parent: root, Depth: 1}
	child2 := &tree.Node{Name: "middle", Parent: root, Depth: 1}
	child3 := &tree.Node{Name: "last", Parent: root, Depth: 1}

	root.Children = []*tree.Node{child1, child2, child3}

	tests := []struct {
		name     string
		node     *tree.Node
		expected bool
	}{
		{"root has no parent", root, true},
		{"first child is not last", child1, false},
		{"middle child is not last", child2, false},
		{"last child is last", child3, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.node.IsLastChild()
			if result != tt.expected {
				t.Errorf("IsLastChild() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTreeBuildInvalidPath(t *testing.T) {
	nonExistentPath := "/this/path/does/not/exist"
	root := tree.Build(nonExistentPath, 1)

	if root == nil {
		t.Fatal("Build should not return nil for invalid path")
	}

	if len(root.Children) != 0 {
		t.Error("Invalid path should result in no children")
	}
}

func TestNodeLoadChildrenInvalidPath(t *testing.T) {
	node := &tree.Node{
		Name:  "invalid",
		Path:  "/invalid/path",
		IsDir: true,
		Depth: 0,
	}

	// Should not panic
	node.LoadChildren()

	if len(node.Children) != 0 {
		t.Error("LoadChildren on invalid path should not add children")
	}
}
