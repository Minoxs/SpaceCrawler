package main

import (
	"log"

	"disk-usage/pkg/DiskExplorer"
)

func main() {
	var dir = DiskExplorer.Map(".")

	log.Println("\n", dir.String())

	for {
		dir.Expand()
		log.Println("Explored")
		log.Println("\n", dir.String())
		if dir.FullyExplored() {
			break
		}
	}
}

// func main() {
// 	var dir = DiskExplorer.Map(".")
//
// 	log.Println("\n", dir.String())
//
// 	var ready = make(chan bool)
//
// 	go func() {
// 		for {
// 			dir.Expand()
// 			log.Println("Exploration done")
// 			ready <- true
// 			if dir.FullyExplored() {
// 				return
// 			}
// 		}
// 	}()
//
// 	for {
// 		// Wait until ready
// 		<-ready
// 		log.Println(dir.Name, dir.Size)
// 		if dir.FullyExplored() {
// 			log.Println("Done")
// 			break
// 		}
// 	}
//
// 	log.SetOutput(os.Stdout)
// 	log.Println("\n", dir.String())
// }
