package main

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/minoxs/SpaceCrawler/pkg/DiskExplorer"
)

func update(node *tview.TreeNode) {
	var info = node.GetReference().(*DiskExplorer.DiskInfo)

	node.SetText(info.String())
	if info.IsDir {
		if info.Explored() {
			node.SetColor(tcell.ColorGreen)
		} else {
			node.SetColor(tcell.ColorOrangeRed)
		}
	} else {
		node.SetColor(tcell.ColorBlue)
	}

	for _, child := range node.GetChildren() {
		update(child)
	}
}

// TODO make travel work recursively without clearing the children every time
func travel(target *tview.TreeNode) {
	target.ClearChildren()
	update(target)

	var info = target.GetReference().(*DiskExplorer.DiskInfo)
	for i, child := range info.Children {
		var node = tview.
			NewTreeNode("").
			SetReference(&info.Children[i]).
			SetSelectable(child.IsDir)

		update(node)
		node.SetExpanded(false)
		target.AddChild(node)
		travel(node)
	}
}

func main() {
	var app = tview.NewApplication()
	var disk = DiskExplorer.Map(".")
	var root = tview.NewTreeNode(disk.Path).SetColor(tcell.ColorRed).SetReference(&disk)
	var tree = tview.NewTreeView().SetRoot(root).SetCurrentNode(root)

	// Add the current directory to the root node.
	travel(root)

	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
			app.QueueUpdateDraw(
				func() {
					update(root)
				},
			)
		}
	}()

	go func() {
		for !disk.Explored() {
			disk.Deepen()
			app.QueueUpdateDraw(
				func() {
					travel(root)
				},
			)
		}
	}()

	// If a directory was selected, open it.
	tree.SetSelectedFunc(
		func(node *tview.TreeNode) {
			var info = node.GetReference().(*DiskExplorer.DiskInfo)
			if info == &disk {
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

	if err := app.SetRoot(tree, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
