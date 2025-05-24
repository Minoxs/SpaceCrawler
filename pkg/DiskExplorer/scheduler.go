package DiskExplorer

type semaphore struct {
	sem chan struct{}
}

func (s *semaphore) Acquire() {
	s.sem <- struct{}{}
}

func (s *semaphore) Release() {
	<-s.sem
}

type Scheduler struct {
	sem semaphore
}

func NewScheduler(max int) (scheduler *Scheduler) {
	scheduler = &Scheduler{
		sem: semaphore{make(chan struct{}, max)},
	}
	return
}

func (s *Scheduler) Go(f func()) {
	s.sem.Acquire()
	go func(f func()) {
		defer s.sem.Release()
		f()
	}(f)
}
