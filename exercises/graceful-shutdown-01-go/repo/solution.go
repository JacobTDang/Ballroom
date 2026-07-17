package main

// Server runs handle on submitted jobs with a pool of workers. Stop
// must drain everything accepted, refuse new work, and only return
// when the workers are done.
//
// TODO: this Stop flips the flag and returns immediately -- queued
// jobs are abandoned and the workers are left running.
type Server struct {
	jobs    chan int
	stopped bool
}

func NewServer(workers int, handle func(int)) *Server {
	s := &Server{jobs: make(chan int, 1024)}
	for w := 0; w < workers; w++ {
		go func() {
			for v := range s.jobs {
				handle(v)
			}
		}()
	}
	return s
}

func (s *Server) Submit(v int) bool {
	if s.stopped {
		return false
	}
	s.jobs <- v
	return true
}

func (s *Server) Stop() {
	s.stopped = true
}
