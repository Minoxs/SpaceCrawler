package DiskExplorer

import "log"

func panicOnError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
