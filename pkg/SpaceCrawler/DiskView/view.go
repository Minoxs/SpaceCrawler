package DiskView

import (
	"context"
	"sync/atomic"
	"time"

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
		disk:     DiskExplorer.Map(root),
		tree:     tview.NewTreeView(),
		syncing:  &atomic.Bool{},
		OnUpdate: func() {},
	}

	var node = tview.NewTreeNode(v.disk.Path).SetReference(&v.disk)
	setNodeInfo(node)

	v.tree.SetBorder(true).SetTitle("Files")
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
	go v.disk.BreadthSearch(context.Background())
	for !v.disk.Explored() {
		expand(v.tree.GetRoot())
		update(v.tree.GetRoot())
		v.OnUpdate()
		time.Sleep(100 * time.Millisecond)
	}
	expand(v.tree.GetRoot())
	update(v.tree.GetRoot())
	v.OnUpdate()
}

// Model returns the model for rendering by tview
func (v *View) Model() tview.Primitive {
	return v.tree
}
