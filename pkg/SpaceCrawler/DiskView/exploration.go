package DiskView

import (
	"github.com/rivo/tview"

	"github.com/minoxs/SpaceCrawler/pkg/DiskExplorer"
)

// expand will recursively expand the folder to find every sub-folder and file in the tree
func expand(target *tview.TreeNode) {
	// target.GetReference().(*DiskExplorer.DiskInfo).Expand()
	travel(target)

	for _, node := range target.GetChildren() {
		if node.GetReference().(*DiskExplorer.DiskInfo).IsDir {
			expand(node)
		}
	}
}

// travel updates a node's children
func travel(target *tview.TreeNode) {
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
