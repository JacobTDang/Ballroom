package main

import "sync"

// Lazy: sync.Once IS this problem, solved -- Do guarantees exactly-once
// execution and makes the result visible to every caller that returns.
type Lazy struct {
	init  func() int
	once  sync.Once
	value int
}

func NewLazy(init func() int) *Lazy {
	return &Lazy{init: init}
}

func (l *Lazy) Get() int {
	l.once.Do(func() {
		l.value = l.init()
	})
	return l.value
}
