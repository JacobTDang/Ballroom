package main

import "testing"

func TestMinStack(t *testing.T) {
	s := NewMinStack()
	s.Push(-2)
	s.Push(0)
	s.Push(-3)
	if got := s.GetMin(); got != -3 {
		t.Fatalf("GetMin() = %d, want -3", got)
	}
	s.Pop()
	if got := s.Top(); got != 0 {
		t.Fatalf("Top() = %d, want 0", got)
	}
	if got := s.GetMin(); got != -2 {
		t.Fatalf("GetMin() = %d, want -2", got)
	}
}

func TestMinStack_MinUpdatesAsEqualValuesArePushedAndPopped(t *testing.T) {
	s := NewMinStack()
	s.Push(1)
	s.Push(1)
	s.Push(1)
	if got := s.GetMin(); got != 1 {
		t.Fatalf("GetMin() = %d, want 1", got)
	}
	s.Pop()
	if got := s.GetMin(); got != 1 {
		t.Fatalf("GetMin() = %d, want 1 (two 1s remain)", got)
	}
	s.Pop()
	s.Pop()
}

func TestMinStack_MinRevertsAfterPoppingNewMin(t *testing.T) {
	s := NewMinStack()
	s.Push(5)
	s.Push(3)
	s.Push(7)
	s.Push(1)
	if got := s.GetMin(); got != 1 {
		t.Fatalf("GetMin() = %d, want 1", got)
	}
	s.Pop() // pop 1
	if got := s.GetMin(); got != 3 {
		t.Fatalf("GetMin() = %d, want 3", got)
	}
}

func TestMinStack_AllSameValuesThenNewMin(t *testing.T) {
	s := NewMinStack()
	s.Push(2)
	s.Push(2)
	s.Push(2)
	if got := s.GetMin(); got != 2 {
		t.Fatalf("GetMin() = %d, want 2", got)
	}
	s.Pop()
	if got := s.GetMin(); got != 2 {
		t.Fatalf("GetMin() = %d, want 2", got)
	}
	s.Pop()
	if got := s.GetMin(); got != 2 {
		t.Fatalf("GetMin() = %d, want 2", got)
	}
	s.Push(-5)
	if got := s.GetMin(); got != -5 {
		t.Fatalf("GetMin() = %d, want -5", got)
	}
}
