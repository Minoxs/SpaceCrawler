package main

import (
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

func travel(target *tview.TreeNode) {
	// update(target)

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

		// update(node)
		// travel(node)
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

	go expand(root)

	app.SetBeforeDrawFunc(
		func(screen tcell.Screen) bool {
			screen.Clear()
			update(root)
			return false
		},
	)

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
