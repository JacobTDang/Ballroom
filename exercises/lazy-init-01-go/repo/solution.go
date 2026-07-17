package main

// Lazy computes a value on first use -- the init function is expensive
// and must run exactly once, no matter how many goroutines call Get
// concurrently.
//
// TODO: the check below isn't atomic with the assignment -- two
// goroutines can both see done == false and both run init.
type Lazy struct {
	init  func() int
	value int
	done  bool
}

func NewLazy(init func() int) *Lazy {
	return &Lazy{init: init}
}

func (l *Lazy) Get() int {
	if !l.done {
		l.value = l.init()
		l.done = true
	}
	return l.value
}
