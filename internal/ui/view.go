package ui

import (
	"dtree/internal/tree"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the TUI display
func (m *Model) View() string {
	var b strings.Builder

	header := m.headerStyle.Render(fmt.Sprintf("DTree - %s (initial depth: %d)", m.rootPath, m.initialDepth))
	b.WriteString(header + "\n\n")

	// Calculate visible range for viewport
	start := m.viewportOffset
	end := start + m.viewportHeight
	if end > len(m.flattenedNodes) {
		end = len(m.flattenedNodes)
	}

	// If viewport can show all nodes, just show everything (for tests and large terminals)
	if m.viewportHeight >= len(m.flattenedNodes) {
		start = 0
		end = len(m.flattenedNodes)
	}

	// Render visible nodes
	for i := start; i < end; i++ {
		node := m.flattenedNodes[i]
		line := m.renderTreeLine(i, node)
		b.WriteString(line + "\n")
	}

	controls := lipgloss.NewStyle().Render("\nControls: ↑↓/jk navigate, Ctrl+U/D half-page, Ctrl+B/F full-page, gg/G top/bottom, Enter/Space expand/collapse, q quit")
	b.WriteString(controls)

	if m.status != "" {
		status := m.errorStyle.Render("\n" + m.status)
		b.WriteString(status)
	}

	return b.String()
}

// renderTreeLine formats a single tree node with styling and tree characters
func (m *Model) renderTreeLine(index int, node *tree.Node) string {
	cursor := " "
	if index == m.cursor {
		cursor = m.cursorStyle.Render(">")
	}

	treeChars := m.getTreeChars(node)

	var nameStyle lipgloss.Style
	var name string

	if node.IsDir {
		nameStyle = m.dirStyle
		var indicator string
		if node.IsExpanded {
			indicator = "▼ "
		} else {
			indicator = "▶ "
		}
		name = nameStyle.Render(indicator + node.Name)
	} else {
		nameStyle = m.fileStyle
		name = nameStyle.Render(node.Name)
	}

	return fmt.Sprintf("%s %s%s", cursor, treeChars, name)
}

// getTreeChars generates proper tree connecting characters (├──, └──, │)
func (m *Model) getTreeChars(node *tree.Node) string {
	if node.Depth == 0 {
		return ""
	}

	var result strings.Builder

	current := node
	positions := make([]bool, node.Depth)

	for current.Parent != nil {
		isLast := current.IsLastChild()
		positions[current.Depth-1] = !isLast
		current = current.Parent
	}

	for i := 0; i < node.Depth-1; i++ {
		if positions[i] {
			result.WriteString("│   ")
		} else {
			result.WriteString("    ")
		}
	}

	if node.IsLastChild() {
		result.WriteString("└── ")
	} else {
		result.WriteString("├── ")
	}

	return result.String()
}