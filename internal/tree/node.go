package tree

import (
	"os"
	"path/filepath"
)

// Node represents a file or directory in the tree structure
type Node struct {
	Name       string
	Path       string
	IsDir      bool
	IsExpanded bool
	Children   []*Node
	Parent     *Node
	Depth      int
}

// Build creates the initial tree structure with specified depth expansion
func Build(rootPath string, initialDepth int) *Node {
	root := &Node{
		Name:       filepath.Base(rootPath),
		Path:       rootPath,
		IsDir:      true,
		IsExpanded: true,
		Depth:      0,
	}

	loadChildrenRecursive(root, initialDepth)
	return root
}

// LoadChildren reads directory contents and creates child nodes
func (n *Node) LoadChildren() {
	entries, err := os.ReadDir(n.Path)
	if err != nil {
		return
	}

	for _, entry := range entries {
		childPath := filepath.Join(n.Path, entry.Name())
		child := &Node{
			Name:   entry.Name(),
			Path:   childPath,
			IsDir:  entry.IsDir(),
			Parent: n,
			Depth:  n.Depth + 1,
		}
		n.Children = append(n.Children, child)
	}
}

// IsLastChild determines if a node is the last child of its parent
func (n *Node) IsLastChild() bool {
	if n.Parent == nil {
		return true
	}

	visibleSiblings := n.Parent.Children

	if len(visibleSiblings) == 0 {
		return true
	}

	return visibleSiblings[len(visibleSiblings)-1] == n
}

// loadChildrenRecursive loads directory contents up to the initial depth
func loadChildrenRecursive(node *Node, initialDepth int) {
	if node.Depth >= initialDepth {
		return
	}

	entries, err := os.ReadDir(node.Path)
	if err != nil {
		return
	}

	for _, entry := range entries {
		childPath := filepath.Join(node.Path, entry.Name())
		child := &Node{
			Name:   entry.Name(),
			Path:   childPath,
			IsDir:  entry.IsDir(),
			Parent: node,
			Depth:  node.Depth + 1,
		}
		node.Children = append(node.Children, child)

		if child.IsDir && child.Depth < initialDepth {
			child.IsExpanded = true
			loadChildrenRecursive(child, initialDepth)
		}
	}
}
