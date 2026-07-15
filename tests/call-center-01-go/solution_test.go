package main

import "testing"

func TestFirstCallGoesToARespondent(t *testing.T) {
	c := NewCallCenter(1, 1, 1)
	if got := c.Dispatch(1); got != "respondent" {
		t.Errorf("Dispatch(1) = %q, want respondent", got)
	}
}

func TestEscalatesToManagerWhenRespondentsBusy(t *testing.T) {
	c := NewCallCenter(1, 1, 1)
	c.Dispatch(1)
	if got := c.Dispatch(2); got != "manager" {
		t.Errorf("Dispatch(2) = %q, want manager", got)
	}
}

func TestEscalatesToDirectorWhenManagersBusy(t *testing.T) {
	c := NewCallCenter(1, 1, 1)
	c.Dispatch(1)
	c.Dispatch(2)
	if got := c.Dispatch(3); got != "director" {
		t.Errorf("Dispatch(3) = %q, want director", got)
	}
}

func TestQueuesWhenEveryoneIsBusy(t *testing.T) {
	c := NewCallCenter(1, 0, 0)
	c.Dispatch(1)
	if got := c.Dispatch(2); got != "queued" {
		t.Errorf("Dispatch(2) = %q, want queued", got)
	}
	if got := c.HandlerOf(2); got != "queued" {
		t.Errorf("HandlerOf(2) = %q, want queued", got)
	}
}

func TestEndingACallAssignsTheQueuedCall(t *testing.T) {
	c := NewCallCenter(1, 1, 0)
	c.Dispatch(1)
	c.Dispatch(2)
	c.Dispatch(3)
	if !c.EndCall(1) {
		t.Fatal("EndCall(1) = false, want true")
	}
	if got := c.HandlerOf(3); got != "respondent" {
		t.Errorf("HandlerOf(3) = %q, want respondent -- the freed employee takes the queued call", got)
	}
}

func TestQueuedCallsAreAssignedFIFO(t *testing.T) {
	c := NewCallCenter(2, 0, 0)
	c.Dispatch(1)
	c.Dispatch(2)
	c.Dispatch(3)
	c.Dispatch(4)
	c.EndCall(2)
	if got := c.HandlerOf(3); got != "respondent" {
		t.Errorf("HandlerOf(3) = %q, want respondent (queued first)", got)
	}
	if got := c.HandlerOf(4); got != "queued" {
		t.Errorf("HandlerOf(4) = %q, want still queued", got)
	}
}

func TestEndCallFreesTheEmployeeWhenQueueIsEmpty(t *testing.T) {
	c := NewCallCenter(1, 0, 0)
	c.Dispatch(1)
	c.EndCall(1)
	if got := c.Dispatch(2); got != "respondent" {
		t.Errorf("Dispatch(2) = %q, want respondent after the slot freed", got)
	}
}

func TestEndUnknownCallReturnsFalse(t *testing.T) {
	if NewCallCenter(1, 1, 1).EndCall(99) {
		t.Error("EndCall(99) = true, want false")
	}
}

func TestEndCallTwiceReturnsFalse(t *testing.T) {
	c := NewCallCenter(1, 0, 0)
	c.Dispatch(1)
	if !c.EndCall(1) {
		t.Fatal("first EndCall = false, want true")
	}
	if c.EndCall(1) {
		t.Error("second EndCall = true, want false")
	}
}

func TestAbandoningAQueuedCallRemovesIt(t *testing.T) {
	c := NewCallCenter(1, 0, 0)
	c.Dispatch(1)
	c.Dispatch(2)
	c.Dispatch(3)
	if !c.EndCall(2) {
		t.Fatal("EndCall(queued 2) = false, want true")
	}
	c.EndCall(1)
	if got := c.HandlerOf(3); got != "respondent" {
		t.Errorf("HandlerOf(3) = %q, want respondent -- abandoned call 2 must be skipped", got)
	}
	if got := c.HandlerOf(2); got != "" {
		t.Errorf("HandlerOf(2) = %q, want empty after abandoning", got)
	}
}

func TestHandlerOfUnknownCallIsEmpty(t *testing.T) {
	if got := NewCallCenter(1, 1, 1).HandlerOf(42); got != "" {
		t.Errorf("HandlerOf(42) = %q, want empty", got)
	}
}

func TestHandlerOfReportsTheActiveLevel(t *testing.T) {
	c := NewCallCenter(1, 1, 1)
	c.Dispatch(1)
	c.Dispatch(2)
	if got := c.HandlerOf(1); got != "respondent" {
		t.Errorf("HandlerOf(1) = %q, want respondent", got)
	}
	if got := c.HandlerOf(2); got != "manager" {
		t.Errorf("HandlerOf(2) = %q, want manager", got)
	}
}
