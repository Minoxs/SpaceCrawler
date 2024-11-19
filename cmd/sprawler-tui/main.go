package main

import (
	"github.com/rivo/tview"

	"github.com/minoxs/SpaceCrawler/pkg/SpaceCrawler"
	"github.com/minoxs/SpaceCrawler/pkg/SpaceCrawler/DiskView"
)

func main() {
	var app = tview.NewApplication()
	app.EnableMouse(true).EnablePaste(false)

	var root = tview.NewFlex().SetDirection(tview.FlexRow)
	tview.NewGrid()
	app.SetRoot(root, true)

	var disk = DiskView.New(".")
	root.AddItem(disk.Model(), 0, 1, true)

	disk.OnUpdate = func() {
		app.QueueUpdateDraw(func() {})
	}
	go disk.Expand()

	var keys = SpaceCrawler.NewKeyMap().
		Add(
			SpaceCrawler.
				NewKeyBind("Quit", []rune{'q', 'Q'}).WithDescription("Quits Program"),
		).
		Add(
			SpaceCrawler.
				NewKeyBind("Other", []rune{'o'}).WithDescription("Goes Bazinga"),
		)
	var view = SpaceCrawler.NewKeyView(keys)
	root.AddItem(view.Model(), 2, 1, false)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
