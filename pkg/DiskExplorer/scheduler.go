package DiskExplorer

import (
	"time"
)

type semaphore struct {
	forks chan struct{}
}

func newSemaphore(max int) *semaphore {
	return &semaphore{
		forks: make(chan struct{}, max),
	}
}

func (s *semaphore) Acquire() {
	s.forks <- struct{}{}
}

func (s *semaphore) Release() {
	<-s.forks
}

func (s *semaphore) Wait() {
	for len(s.forks) > 0 {
		time.Sleep(100 * time.Millisecond)
	}
}
