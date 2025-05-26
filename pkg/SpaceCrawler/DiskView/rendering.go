package DiskView

import (
	"context"
	"sort"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/minoxs/SpaceCrawler/pkg/DiskExplorer"
)

// synchronize will keep the model updated as new files are discovered
func (v *View) synchronize(ctx context.Context) {
	for ctx.Err() == nil {
		time.Sleep(100 * time.Millisecond)
		update(v.tree.GetRoot())
		v.OnUpdate()
	}
	update(v.tree.GetRoot())
	v.OnUpdate()
}

// update will update node information
func update(node *tview.TreeNode) {
	setNodeInfo(node)
	if !node.IsExpanded() {
		return
	}

	setChildren(node)
	for _, child := range node.GetChildren() {
		update(child)
	}
	sortNodeInfo(node)
}

// setChildren updates a node's children
func setChildren(target *tview.TreeNode) {
	var info = target.GetReference().(*DiskExplorer.DiskInfo)
	if len(info.Children) == len(target.GetChildren()) {
		return
	}

	target.ClearChildren()
	for i, child := range info.Children {
		var node = tview.
			NewTreeNode("").
			SetReference(&info.Children[i]).
			SetSelectable(child.IsDir)

		node.SetExpanded(false)
		target.AddChild(node)
	}
}

// setNodeInfo will add the information to the node for rendering
func setNodeInfo(node *tview.TreeNode) {
	var info = node.GetReference().(*DiskExplorer.DiskInfo)

	// Set node text
	node.SetText(info.String())

	// Color switch
	switch {
	case info.Denied():
		node.SetColor(tcell.ColorOrange)
	case !info.IsDir:
		node.SetColor(tcell.ColorBlue)
	case info.Explored():
		node.SetColor(tcell.ColorGreen)
	default:
		node.SetColor(tcell.ColorRed)
	}
}

// sortNodeInfo will sort the node information
func sortNodeInfo(node *tview.TreeNode) {
	var children = node.GetChildren()
	sort.Slice(
		children, func(i, j int) bool {
			var a = children[i].GetReference().(*DiskExplorer.DiskInfo)
			var b = children[j].GetReference().(*DiskExplorer.DiskInfo)
			return a.Size() > b.Size()
		},
	)
	node.SetChildren(children)
}
