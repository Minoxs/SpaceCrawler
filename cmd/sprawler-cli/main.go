package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/minoxs/SpaceCrawler/pkg/DiskExplorer"
)

func main() {
	var dir = DiskExplorer.Map(".")
	log.Println("\n", dir.String())

	var ctx, cancelFunc = context.WithCancel(context.Background())
	var ready = make(chan bool)
	var interrupt = make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Starting Walk")
		dir.Exhaust(ctx)
		ready <- true
	}()

	select {
	case <-ready:
		log.Println("Ready")
	case <-interrupt:
		log.Println("Canceled")
		cancelFunc()
	}

	if dir.Explored() {
		log.Println("Fully explored")
	} else {
		log.Println("Partial analysis")
	}

	dir.Render(os.Stdout)
}
