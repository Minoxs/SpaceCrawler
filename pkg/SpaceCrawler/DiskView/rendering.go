package DiskView

import (
	"sort"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/minoxs/SpaceCrawler/pkg/DiskExplorer"
)

// synchronize will keep the model updated as new files are discovered
func (v *View) synchronize() {
	if v.syncing.CompareAndSwap(false, true) {
		for {
			var explored = v.disk.Explored()
			update(v.tree.GetRoot())
			v.OnUpdate()
			if explored {
				break
			}
			time.Sleep(250 * time.Millisecond)
		}

		v.syncing.Store(false)
	}
}

// update will update node information
func update(node *tview.TreeNode) {
	setNodeInfo(node)
	for _, child := range node.GetChildren() {
		update(child)
	}
	sortNodeInfo(node)
}

// setNodeInfo will add the information to the node for rendering
func setNodeInfo(node *tview.TreeNode) {
	var info = node.GetReference().(*DiskExplorer.DiskInfo)

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
