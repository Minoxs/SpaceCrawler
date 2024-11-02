package main

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/minoxs/SpaceCrawler/pkg/DiskExplorer"
)

func setNodeInfo(node *tview.TreeNode) {
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
}

func update(node *tview.TreeNode) {
	setNodeInfo(node)
	for _, child := range node.GetChildren() {
		update(child)
	}
}

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

func expand(target *tview.TreeNode) {
	target.GetReference().(*DiskExplorer.DiskInfo).Expand()
	travel(target)

	for _, node := range target.GetChildren() {
		if node.GetReference().(*DiskExplorer.DiskInfo).IsDir {
			go expand(node)
		}
	}
}

func main() {
	var app = tview.NewApplication()
	var disk = DiskExplorer.Map(".")
	var root = tview.NewTreeNode(disk.Path).SetColor(tcell.ColorRed).SetReference(&disk)
	var tree = tview.NewTreeView().SetRoot(root).SetCurrentNode(root)

	setNodeInfo(root)
	go expand(root)

	go func() {
		for {
			time.Sleep(100 * time.Millisecond)

			var explored = disk.Explored()
			app.QueueUpdateDraw(
				func() {
					update(root)
				},
			)

			if explored {
				return
			}
		}
	}()

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
