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
	app.EnableMouse(true).EnablePaste(false)

	var disk = DiskExplorer.Map(".")

	// var header = tview.NewFlex().SetDirection(tview.FlexColumn).
	// 	AddItem(
	// 		tview.NewTextView().SetText("Search: "), 8, 1, false,
	// 	).
	// 	AddItem(
	// 		tview.NewTextArea().SetText("HERE", true), 0, 1, true,
	// 	)
	//
	// var text = tview.NewTextView().
	// 	SetDynamicColors(true).
	// 	SetRegions(true).
	// 	SetWrap(false)

	var root = tview.NewTreeNode(disk.Path).SetColor(tcell.ColorRed).SetReference(&disk)
	var tree = tview.NewTreeView().SetRoot(root).SetCurrentNode(root)
	tree.SetBorder(true).SetTitle("Files")
	app.SetRoot(tree, true)

	// var box = tview.NewFlex().
	// 	SetDirection(tview.FlexRow).
	// 	AddItem(header, 1, 1, false).
	// 	AddItem(tree, 0, 1, true).
	// 	AddItem(text, 1, 1, false)
	// app.SetRoot(box, true)

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

	app.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Rune() {
			case 'q', 'Q':
				app.Stop()
				return nil
			}

			return event
		},
	)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
