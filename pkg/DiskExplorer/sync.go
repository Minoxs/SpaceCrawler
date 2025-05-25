package DiskExplorer

import (
	"time"
)

// semaphore ensures only N threads can access some resource
type semaphore struct {
	forks chan struct{}
}

// newSemaphore creates a semaphore with max resources
func newSemaphore(max int) *semaphore {
	return &semaphore{
		forks: make(chan struct{}, max),
	}
}

// Acquire enters the semaphore
func (s *semaphore) Acquire() {
	s.forks <- struct{}{}
}

// Release releases the semaphore
func (s *semaphore) Release() {
	<-s.forks
}

// Wait will block until the semaphore count reaches 0
func (s *semaphore) Wait() {
	for len(s.forks) > 0 {
		time.Sleep(100 * time.Millisecond)
	}
}
