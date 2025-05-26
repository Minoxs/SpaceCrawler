package main

import (
	"context"

	"github.com/rivo/tview"

	"github.com/minoxs/SpaceCrawler/pkg/SpaceCrawler/DiskView"
)

func NoOp() {

}

func main() {
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	var app = tview.NewApplication()
	app.EnableMouse(true).EnablePaste(false)

	var root = DiskView.New(".")
	app.SetRoot(root.Model(), true)
	root.OnUpdate = func() {
		app.QueueUpdateDraw(NoOp)
	}

	go root.Expand(ctx)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
