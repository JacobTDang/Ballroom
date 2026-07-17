package main

import (
	"reflect"
	"testing"
)

func TestRegistrationOrderAndEventIsolation(t *testing.T) {
	e := NewEmitter()
	var got []string
	e.On("a", func(v int) { got = append(got, "first") })
	e.On("a", func(v int) { got = append(got, "second") })
	e.On("b", func(v int) { got = append(got, "other-event") })

	e.Emit("a", 1)
	want := []string{"first", "second"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("emit order = %v, want %v (registration order, own event only)", got, want)
	}
}

func TestOnceFiresExactlyOnce(t *testing.T) {
	e := NewEmitter()
	calls := 0
	e.Once("a", func(v int) { calls++ })
	e.Emit("a", 1)
	e.Emit("a", 2)
	e.Emit("a", 3)
	if calls != 1 {
		t.Fatalf("once handler fired %d times, want exactly 1", calls)
	}
}

func TestOffRemovesExactlyThatSubscription(t *testing.T) {
	e := NewEmitter()
	var got []string
	id := e.On("a", func(v int) { got = append(got, "removed") })
	e.On("a", func(v int) { got = append(got, "kept") })
	e.Off(id)
	e.Emit("a", 1)
	if !reflect.DeepEqual(got, []string{"kept"}) {
		t.Fatalf("after Off, emit called %v, want only the kept handler", got)
	}
}

func TestRemovalDuringEmitSkipsTheRemoved(t *testing.T) {
	e := NewEmitter()
	var got []string
	var victim int
	e.On("a", func(v int) {
		got = append(got, "assassin")
		e.Off(victim)
	})
	victim = e.On("a", func(v int) { got = append(got, "victim") })
	e.On("a", func(v int) { got = append(got, "bystander") })

	e.Emit("a", 1)
	want := []string{"assassin", "bystander"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("emit with mid-flight removal = %v, want %v (the removed handler must not fire)", got, want)
	}
}

func TestEmitWithNoHandlersIsANoOp(t *testing.T) {
	e := NewEmitter()
	e.Emit("nobody-listening", 42) // must simply not panic
}

func TestHandlersReceiveTheValue(t *testing.T) {
	e := NewEmitter()
	var got int
	e.On("a", func(v int) { got = v })
	e.Emit("a", 99)
	if got != 99 {
		t.Fatalf("handler received %d, want 99", got)
	}
}
