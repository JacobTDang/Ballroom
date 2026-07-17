package main

import "sync"

// Server: the closed jobs channel is the "no more work" signal -- the
// workers' range loops end when it's drained, and the WaitGroup is how
// Stop waits for that. The mutex makes Submit's stopped-check atomic
// with its send, so no Submit can sneak a job into a closing channel.
type Server struct {
	mu      sync.Mutex
	jobs    chan int
	stopped bool
	wg      sync.WaitGroup
}

func NewServer(workers int, handle func(int)) *Server {
	s := &Server{jobs: make(chan int, 1024)}
	for w := 0; w < workers; w++ {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			for v := range s.jobs {
				handle(v)
			}
		}()
	}
	return s
}

func (s *Server) Submit(v int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.stopped {
		return false
	}
	s.jobs <- v
	return true
}

func (s *Server) Stop() {
	s.mu.Lock()
	if s.stopped {
		s.mu.Unlock()
		s.wg.Wait()
		return
	}
	s.stopped = true
	close(s.jobs)
	s.mu.Unlock()
	s.wg.Wait()
}
