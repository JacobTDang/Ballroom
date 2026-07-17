package main

// Emitter: On/Once subscribe (returning an id), Off unsubscribes,
// Emit calls the event's handlers in registration order.
//
// TODO: no ids (always 0), Off does nothing, and Once is just On --
// it never unhooks itself.
type Emitter struct {
	handlers map[string][]func(int)
}

func NewEmitter() *Emitter {
	return &Emitter{handlers: make(map[string][]func(int))}
}

func (e *Emitter) On(event string, fn func(int)) int {
	e.handlers[event] = append(e.handlers[event], fn)
	return 0
}

func (e *Emitter) Once(event string, fn func(int)) int {
	return e.On(event, fn)
}

func (e *Emitter) Off(id int) {}

func (e *Emitter) Emit(event string, v int) {
	for _, fn := range e.handlers[event] {
		fn(v)
	}
}
