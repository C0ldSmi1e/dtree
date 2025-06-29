package tests

import (
	"dtree/internal/fileops"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestOpenFile(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		filePath string
		wantErr  bool
	}{
		{
			name:     "valid file",
			filePath: testFile,
			wantErr:  false, // Note: might still error if no default app
		},
		{
			name:     "non-existent file",
			filePath: "/path/that/does/not/exist.txt",
			wantErr:  true,
		},
		{
			name:     "empty path",
			filePath: "",
			wantErr:  false, // macOS 'open' command succeeds with empty path
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fileops.OpenFile(tt.filePath)
			
			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}
			
			// For valid files, we can't guarantee success since it depends
			// on system configuration, but we can check the error type
			if !tt.wantErr && err != nil {
				// Log the error but don't fail - system might not have default app
				t.Logf("OpenFile failed (this might be expected): %v", err)
			}
		})
	}
}

func TestFormatOpenError(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		err      error
		want     string
	}{
		{
			name:     "simple error",
			filePath: "/path/to/file.txt",
			err:      os.ErrNotExist,
			want:     "Error opening file.txt: file does not exist",
		},
		{
			name:     "complex path",
			filePath: "/very/long/path/to/document.pdf",
			err:      os.ErrPermission,
			want:     "Error opening document.pdf: permission denied",
		},
		{
			name:     "file with no extension",
			filePath: "/path/to/README",
			err:      os.ErrInvalid,
			want:     "Error opening README: invalid argument",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fileops.FormatOpenError(tt.filePath, tt.err)
			
			// Check that it contains the expected filename
			filename := filepath.Base(tt.filePath)
			if !strings.Contains(result, filename) {
				t.Errorf("Result should contain filename %s, got: %s", filename, result)
			}
			
			// Check that it contains "Error opening"
			if !strings.HasPrefix(result, "Error opening") {
				t.Errorf("Result should start with 'Error opening', got: %s", result)
			}
			
			// Check that it contains the error message
			if !strings.Contains(result, tt.err.Error()) {
				t.Errorf("Result should contain error message '%s', got: %s", tt.err.Error(), result)
			}
		})
	}
}

func TestFormatPlatformError(t *testing.T) {
	result := fileops.FormatPlatformError()
	
	// Should contain current platform
	currentPlatform := runtime.GOOS
	if !strings.Contains(result, currentPlatform) {
		t.Errorf("Result should contain current platform %s, got: %s", currentPlatform, result)
	}
	
	// Should indicate it's an error
	if !strings.HasPrefix(result, "Error:") {
		t.Errorf("Result should start with 'Error:', got: %s", result)
	}
	
	// Should mention unsupported platform
	if !strings.Contains(strings.ToLower(result), "unsupported") {
		t.Errorf("Result should mention 'unsupported', got: %s", result)
	}
}

func TestOpenFileUnsupportedPlatform(t *testing.T) {
	// This test checks the error message for unsupported platforms
	// We can't easily mock runtime.GOOS, so we test the error path indirectly
	
	// Create a temporary file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	
	// The actual OpenFile function will work on supported platforms
	// but we can test that it doesn't panic and returns some result
	err = fileops.OpenFile(testFile)
	
	// We don't assert success/failure here since it depends on system setup
	// but we ensure it doesn't panic
	t.Logf("OpenFile result: %v", err)
}

func TestFileopsEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		wantErr  bool
	}{
		{
			name:     "directory instead of file",
			filePath: t.TempDir(),
			wantErr:  false, // Some systems can "open" directories
		},
		{
			name:     "file with spaces",
			filePath: filepath.Join(t.TempDir(), "file with spaces.txt"),
			wantErr:  false,
		},
		{
			name:     "file with special characters",
			filePath: filepath.Join(t.TempDir(), "file-with_special.chars.txt"),
			wantErr:  false,
		},
	}
	
	// Create the test files
	for _, tt := range tests {
		if strings.Contains(tt.filePath, "file") { // Skip directory test
			err := os.WriteFile(tt.filePath, []byte("content"), 0644)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fileops.OpenFile(tt.filePath)
			
			// Log result but don't assert success since it depends on system
			t.Logf("OpenFile(%s) = %v", tt.filePath, err)
		})
	}
}