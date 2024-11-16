package main

import (
	"github.com/rivo/tview"

	"github.com/minoxs/SpaceCrawler/pkg/SpaceCrawler/DiskView"
)

func main() {
	var app = tview.NewApplication()
	app.EnableMouse(true).EnablePaste(false)

	var root = DiskView.New(".")
	app.SetRoot(root.Model(), true)
	root.OnUpdate = func() {
		app.QueueUpdateDraw(func() {})
	}

	go root.Expand()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
