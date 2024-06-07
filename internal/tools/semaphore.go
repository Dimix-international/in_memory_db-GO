package tools

import "sync"

type Semaphore struct {
	count    int
	maxCount int
	cond     *sync.Cond
}

func NewSemaphore(maxConnections int) *Semaphore {
	mu := &sync.Mutex{}
	return &Semaphore{
		count:    0,
		maxCount: maxConnections,
		cond:     sync.NewCond(mu),
	}
}

func (s *Semaphore) Acquire() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	if s.count >= s.maxCount {
		s.cond.Wait()
	}

	s.count++
}

func (s *Semaphore) Release() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	s.count--
	s.cond.Signal()
}

func (s *Semaphore) WithSemaphore(action func()) {
	if action == nil {
		return
	}
	s.Acquire()
	action()
	s.Release()
}
