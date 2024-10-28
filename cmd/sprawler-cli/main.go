package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/minoxs/SpaceCrawler/pkg/DiskExplorer"
)

func main() {
	var dir = DiskExplorer.Map(".")
	log.Println("\n", dir.String())

	var ready = make(chan bool)
	var cancel = make(chan os.Signal, 1)
	signal.Notify(cancel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		time.Sleep(500 * time.Millisecond)
		cancel <- syscall.SIGINT
	}()

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
				if dir.IsExplored() {
					ready <- true
					return
				}
			}
		}
	}()

	<-ready
	if dir.IsExplored() {
		log.Println("Fully explored")
	} else {
		log.Println("Partial analysis")
	}

	log.SetOutput(os.Stdout)
	log.Println("\n", dir.Render())
}
