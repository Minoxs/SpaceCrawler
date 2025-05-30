package DiskView

import (
	"context"
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

			var expanded = node.IsExpanded()
			node.SetExpanded(!expanded)
			update(node)
		},
	)

	return v
}

// Expand will keep the view updated
func (v *View) Expand(ctx context.Context) {
	if !v.syncing.CompareAndSwap(false, true) {
		return
	}
	defer v.syncing.Store(false)

	var inner, cancel = context.WithCancel(ctx)
	defer cancel()

	go v.synchronize(inner)
	v.disk.BreadthSearch(inner)
}

// Model returns the model for rendering by tview
func (v *View) Model() tview.Primitive {
	return v.tree
}
