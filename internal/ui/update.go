package ui

import (
	"dtree/internal/fileops"

	tea "github.com/charmbracelet/bubbletea"
)

// Update handles keyboard input and state changes
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Update terminal dimensions and recalculate viewport
		m.terminalHeight = msg.Height
		m.terminalWidth = msg.Width
		m.updateViewportHeight()
		m.adjustViewportToCursor()
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				m.adjustViewportToCursor()
			}
		case "down", "j":
			if m.cursor < len(m.flattenedNodes)-1 {
				m.cursor++
				m.adjustViewportToCursor()
			}
		case "ctrl+u":
			m.jumpHalfScreen(-1) // Up
		case "ctrl+d":
			m.jumpHalfScreen(1) // Down
		case "ctrl+b":
			m.jumpFullScreen(-1) // Up
		case "ctrl+f":
			m.jumpFullScreen(1) // Down
		case "g":
			if m.pendingG {
				// Second 'g' - go to top
				m.jumpToTop()
				m.pendingG = false
			} else {
				// First 'g' - wait for second
				m.pendingG = true
			}
		case "G":
			m.jumpToBottom()
			m.pendingG = false // Reset if was waiting for gg
		case "enter", " ":
			m.pendingG = false // Reset pending g on other actions
			if m.cursor < len(m.flattenedNodes) {
				node := m.flattenedNodes[m.cursor]
				if node.IsDir {
					node.IsExpanded = !node.IsExpanded
					// Lazy load children when expanding
					if node.IsExpanded && len(node.Children) == 0 {
						node.LoadChildren()
					}
					m.updateFlattenedNodes()
					m.adjustViewportToCursor()
				} else {
					// Open file with default application
					m.openFileWithStatus(node.Path)
				}
			}
		default:
			// Reset pending 'g' on any other key
			m.pendingG = false
		}
	}
	return m, nil
}

// openFileWithStatus opens a file and updates status accordingly
func (m *Model) openFileWithStatus(filePath string) {
	err := fileops.OpenFile(filePath)
	if err != nil {
		m.status = fileops.FormatOpenError(filePath, err)
	} else {
		m.status = ""
	}
}