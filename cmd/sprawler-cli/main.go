package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/minoxs/SpaceCrawler/pkg/DiskExplorer"
)

func main() {
	var dir = DiskExplorer.Map(".")
	log.Println("\n", dir.String())

	var ready = make(chan bool)
	var cancel = make(chan os.Signal, 1)
	signal.Notify(cancel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		var i = 0
		for {
			i += 1
			dir.Deepen()
			log.Println("Exploration done", i)
			select {
			case <-cancel:
				ready <- true
				return
			default:
				if dir.Explored() {
					ready <- true
					return
				}
			}
		}
	}()

	<-ready
	if dir.Explored() {
		log.Println("Fully explored")
	} else {
		log.Println("Partial analysis")
	}

	log.SetOutput(os.Stdout)
	log.Println("\n", dir.Render())
}
