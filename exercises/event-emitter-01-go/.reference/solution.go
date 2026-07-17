package main

// Emitter: every subscription gets a unique id and lives in a single
// ordered slice per event. Emit walks a snapshot of ids (not
// functions) and re-checks liveness before each call -- that's what
// makes removal-during-emit safe: an id Off'd by an earlier handler
// fails the liveness check and is skipped. Once is just On plus a
// self-removing wrapper.
type subscription struct {
	id   int
	fn   func(int)
	once bool
}

type Emitter struct {
	subs   map[string][]*subscription
	byID   map[int]string // id -> event, for Off
	nextID int
}

func NewEmitter() *Emitter {
	return &Emitter{subs: make(map[string][]*subscription), byID: make(map[int]string)}
}

func (e *Emitter) add(event string, fn func(int), once bool) int {
	e.nextID++
	id := e.nextID
	e.subs[event] = append(e.subs[event], &subscription{id: id, fn: fn, once: once})
	e.byID[id] = event
	return id
}

func (e *Emitter) On(event string, fn func(int)) int {
	return e.add(event, fn, false)
}

func (e *Emitter) Once(event string, fn func(int)) int {
	return e.add(event, fn, true)
}

func (e *Emitter) Off(id int) {
	event, ok := e.byID[id]
	if !ok {
		return
	}
	delete(e.byID, id)
	list := e.subs[event]
	for i, s := range list {
		if s.id == id {
			e.subs[event] = append(list[:i], list[i+1:]...)
			return
		}
	}
}

func (e *Emitter) live(id int) *subscription {
	event, ok := e.byID[id]
	if !ok {
		return nil
	}
	for _, s := range e.subs[event] {
		if s.id == id {
			return s
		}
	}
	return nil
}

func (e *Emitter) Emit(event string, v int) {
	// Snapshot the ids: handlers may subscribe/unsubscribe during the
	// emit, and iterating the live slice while it mutates would skip
	// or double-call neighbors.
	ids := make([]int, 0, len(e.subs[event]))
	for _, s := range e.subs[event] {
		ids = append(ids, s.id)
	}
	for _, id := range ids {
		s := e.live(id)
		if s == nil {
			continue // removed during this emit
		}
		if s.once {
			e.Off(s.id) // unhook BEFORE calling: re-entrant emits must not double-fire
		}
		s.fn(v)
	}
}
