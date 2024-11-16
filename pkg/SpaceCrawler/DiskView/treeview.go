package DiskView

import (
	"sync/atomic"
	"time"

	"github.com/gdamore/tcell/v2"
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

// expand will recursively expand the folder to find every sub-folder and file in the tree
func expand(target *tview.TreeNode) {
	target.GetReference().(*DiskExplorer.DiskInfo).Expand()
	travel(target)

	for _, node := range target.GetChildren() {
		if node.GetReference().(*DiskExplorer.DiskInfo).IsDir {
			go expand(node)
		}
	}
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

// update will update node information
func update(node *tview.TreeNode) {
	setNodeInfo(node)
	for _, child := range node.GetChildren() {
		update(child)
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
