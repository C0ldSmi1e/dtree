package fileops

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
)

// OpenFile opens a file with the default system application
func OpenFile(filePath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", filePath)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", filePath)
	case "linux":
		cmd = exec.Command("xdg-open", filePath)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Run()
}

// FormatOpenError formats an error message for file opening failures
func FormatOpenError(filePath string, err error) string {
	return fmt.Sprintf("Error opening %s: %v", filepath.Base(filePath), err)
}

// FormatPlatformError formats an error message for unsupported platforms
func FormatPlatformError() string {
	return fmt.Sprintf("Error: Unsupported platform %s", runtime.GOOS)
}
