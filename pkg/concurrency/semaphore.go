package concurrency

type token struct{}

type Semaphore struct {
	queue chan token
}

func NewSemaphore(max int) Semaphore {
	return Semaphore{
		queue: make(chan token, max),
	}
}

func (s *Semaphore) Acquire() {
	s.queue <- token{}
}

func (s *Semaphore) Release() {
	<-s.queue
}

func (s *Semaphore) WithSemaphore(f func()) {
	if f == nil {
		return
	}
	s.Acquire()
	defer s.Release()
	f()
}
