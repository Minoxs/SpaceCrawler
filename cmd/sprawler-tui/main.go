package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/minoxs/SpaceCrawler/pkg/DiskExplorer"
)

func main() {
	var app = tview.NewApplication()
	var disk = DiskExplorer.Map(".")
	var root = tview.NewTreeNode(disk.Path).SetColor(tcell.ColorRed).SetReference(&disk)
	var tree = tview.NewTreeView().SetRoot(root).SetCurrentNode(root)

	travel := func(target *tview.TreeNode) {
		var info = target.GetReference().(*DiskExplorer.DiskInfo)

		for i, child := range info.Children {
			var node = tview.
				NewTreeNode(child.String()).
				SetReference(&info.Children[i]).
				SetSelectable(child.IsDir)

			if child.IsDir {
				if child.IsExplored {
					node.SetColor(tcell.ColorGreen)
				} else {
					node.SetColor(tcell.ColorOrangeRed)
				}
			} else {
				node.SetColor(tcell.ColorBlue)
			}

			target.AddChild(node)
		}
	}

	// Add the current directory to the root node.
	travel(root)

	// If a directory was selected, open it.
	tree.SetSelectedFunc(
		func(node *tview.TreeNode) {
			var info = node.GetReference().(*DiskExplorer.DiskInfo)
			if info == &disk {
				return
			}

			var children = node.GetChildren()
			if len(children) == 0 {
				info.Expand()
				travel(node)
				node.SetExpanded(true)
			} else {
				// Collapse if visible, expand if collapsed.
				node.SetExpanded(!node.IsExpanded())
			}
		},
	)

	if err := app.SetRoot(tree, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
