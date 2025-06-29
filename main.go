package main

import (
	"dtree/internal/tree"
	"dtree/internal/ui"
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var initialDepth int
	var rootPath string
	var showHelp bool

	flag.IntVar(&initialDepth, "d", 1, "Initial depth to expand")
	flag.IntVar(&initialDepth, "depth", 1, "Initial depth to expand")
	flag.BoolVar(&showHelp, "h", false, "Show help message")
	flag.BoolVar(&showHelp, "help", false, "Show help message")
	flag.Parse()

	if showHelp {
		fmt.Println("DTree - Interactive directory tree viewer")
		fmt.Println("\nUsage:")
		fmt.Println("  dtree [options] [directory]")
		fmt.Println("\nOptions:")
		fmt.Println("  -d, --depth <num>   Initial depth to expand (default: 1)")
		fmt.Println("  -h, --help          Show this help message")
		fmt.Println("\nControls:")
		fmt.Println("  ↑/↓ or j/k          Navigate up/down")
		fmt.Println("  Ctrl+U/D            Jump half-screen up/down")
		fmt.Println("  Ctrl+B/F            Jump full-screen up/down")
		fmt.Println("  gg/G                Go to top/bottom")
		fmt.Println("  Enter/Space         Expand/collapse directories")
		fmt.Println("  Enter               Open files with default application")
		fmt.Println("  q/Ctrl+C/Esc        Quit")
		fmt.Println("\nExamples:")
		fmt.Println("  dtree               # View current directory")
		fmt.Println("  dtree /home/user    # View specific directory")
		fmt.Println("  dtree -d 3 .        # Expand 3 levels deep")
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) > 0 {
		rootPath = args[0]
	} else {
		var err error
		rootPath, err = os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
	}

	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		fmt.Printf("Directory does not exist: %s\n", rootPath)
		os.Exit(1)
	}

	// Build the tree structure
	rootTree := tree.Build(rootPath, initialDepth)

	// Create the UI model
	model := ui.New(rootTree, initialDepth, rootPath)

	// Run the TUI
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
