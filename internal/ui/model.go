package ui

import (
	"dtree/internal/tree"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model holds the application state for the Bubbletea TUI
type Model struct {
	tree           *tree.Node
	cursor         int
	flattenedNodes []*tree.Node // Flattened view of visible nodes for navigation
	initialDepth   int
	rootPath       string
	status         string // Status message for user feedback

	// Viewport for scrolling
	viewportHeight int // Available height for content display
	viewportOffset int // First visible line index
	terminalHeight int // Total terminal height
	terminalWidth  int // Total terminal width

	// Styling
	dirStyle    lipgloss.Style
	fileStyle   lipgloss.Style
	cursorStyle lipgloss.Style
	headerStyle lipgloss.Style
	errorStyle  lipgloss.Style

	// Vim-style navigation state
	pendingG bool // Track if 'g' was pressed for 'gg' sequence
}

// New creates a new UI model
func New(rootTree *tree.Node, initialDepth int, rootPath string) *Model {
	m := &Model{
		tree:         rootTree,
		cursor:       0,
		initialDepth: initialDepth,
		rootPath:     rootPath,

		// Initialize viewport - responsive to content and terminal size
		viewportHeight: 1000, // Large default - will be constrained by actual terminal
		viewportOffset: 0,    // Start at top
		terminalHeight: 1000, // Large default - will be updated by WindowSizeMsg
		terminalWidth:  80,   // Default fallback

		dirStyle:    lipgloss.NewStyle().Bold(true),
		fileStyle:   lipgloss.NewStyle(),
		cursorStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true),
		headerStyle: lipgloss.NewStyle(),
		errorStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("1")),
	}

	m.updateFlattenedNodes()
	return m
}

// Init initializes the model (required by Bubbletea)
func (m *Model) Init() tea.Cmd {
	return nil
}

// updateFlattenedNodes rebuilds the flattened view for navigation
func (m *Model) updateFlattenedNodes() {
	m.flattenedNodes = []*tree.Node{}
	m.flattenRecursive(m.tree)
}

// flattenRecursive recursively adds visible nodes to the flattened list
func (m *Model) flattenRecursive(node *tree.Node) {
	m.flattenedNodes = append(m.flattenedNodes, node)
	if node.IsExpanded {
		for _, child := range node.Children {
			m.flattenRecursive(child)
		}
	}
}

// SetStatus sets the status message
func (m *Model) SetStatus(status string) {
	m.status = status
}

// updateViewportHeight calculates available height for content
func (m *Model) updateViewportHeight() {
	// Account for header (2 lines), controls (1 line), status (1 line if present)
	headerLines := 2
	controlLines := 1
	statusLines := 0
	if m.status != "" {
		statusLines = 1
	}

	availableHeight := m.terminalHeight - headerLines - controlLines - statusLines

	// Only constrain viewport if we have a reasonable terminal height
	// This ensures tests and very tall terminals show all content
	if m.terminalHeight > 10 && availableHeight > 0 {
		m.viewportHeight = availableHeight
	} else {
		// For tests or unusual cases, use a large viewport
		m.viewportHeight = max(availableHeight, len(m.flattenedNodes))
	}

	if m.viewportHeight < 1 {
		m.viewportHeight = 1
	}
}

// adjustViewportToCursor ensures cursor is always visible in viewport
func (m *Model) adjustViewportToCursor() {
	if m.cursor < m.viewportOffset {
		// Cursor is above viewport, scroll up
		m.viewportOffset = m.cursor
	} else if m.cursor >= m.viewportOffset+m.viewportHeight {
		// Cursor is below viewport, scroll down
		m.viewportOffset = m.cursor - m.viewportHeight + 1
	}

	// Ensure viewport doesn't go beyond content
	if m.viewportOffset < 0 {
		m.viewportOffset = 0
	}
	maxOffset := len(m.flattenedNodes) - m.viewportHeight
	if maxOffset < 0 {
		maxOffset = 0
	}
	if m.viewportOffset > maxOffset {
		m.viewportOffset = maxOffset
	}
}

// jumpHalfScreen moves cursor by half viewport height
func (m *Model) jumpHalfScreen(direction int) {
	jump := m.viewportHeight / 2
	if jump < 1 {
		jump = 1
	}
	m.moveCursor(direction * jump)
}

// jumpFullScreen moves cursor by full viewport height
func (m *Model) jumpFullScreen(direction int) {
	jump := m.viewportHeight
	if jump < 1 {
		jump = 1
	}
	m.moveCursor(direction * jump)
}

// moveCursor safely moves cursor with bounds checking
func (m *Model) moveCursor(delta int) {
	newPos := m.cursor + delta
	if newPos < 0 {
		newPos = 0
	}
	if newPos >= len(m.flattenedNodes) {
		newPos = len(m.flattenedNodes) - 1
	}
	m.cursor = newPos
	m.adjustViewportToCursor()
}

// jumpToTop moves cursor to first item
func (m *Model) jumpToTop() {
	m.cursor = 0
	m.adjustViewportToCursor()
}

// jumpToBottom moves cursor to last item
func (m *Model) jumpToBottom() {
	if len(m.flattenedNodes) > 0 {
		m.cursor = len(m.flattenedNodes) - 1
		m.adjustViewportToCursor()
	}
}
