package DiskView

import (
	"sync/atomic"

	"github.com/rivo/tview"

	"github.com/minoxs/SpaceCrawler/pkg/DiskExplorer"
)

// View is the exported component for the DiskView
// Contains every method required to render out a folder
type View struct {
	disk    DiskExplorer.DiskInfo
	tree    *tview.TreeView
	syncing *atomic.Bool

	// Function called after model updates
	OnUpdate func()
}

// New returns a new instance of the DiskView
func New(root string) (v View) {
	v = View{
		disk:    DiskExplorer.Map(root),
		tree:    tview.NewTreeView(),
		syncing: &atomic.Bool{},
	}

	var node = tview.NewTreeNode(v.disk.Path).SetReference(&v.disk)
	setNodeInfo(node)
	v.tree.SetRoot(node)
	v.tree.SetSelectedFunc(
		func(node *tview.TreeNode) {
			var info = node.GetReference().(*DiskExplorer.DiskInfo)
			if info == &v.disk {
				return
			}

			if info.Expand() {
				travel(node)
				node.SetExpanded(true)
			} else {
				node.SetExpanded(!node.IsExpanded())
			}
		},
	)

	return v
}

// Expand will iterate recursively over the folder
func (v *View) Expand() {
	go v.synchronize()
	expand(v.tree.GetRoot())
}

// Model returns the model for rendering by tview
func (v *View) Model() tview.Primitive {
	return v.tree
}
