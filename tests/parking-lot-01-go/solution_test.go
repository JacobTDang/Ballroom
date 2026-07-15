package main

import "testing"

func TestInitialAvailabilityMatchesConstruction(t *testing.T) {
	lot := NewParkingLot(2, 3, 4)
	if got := lot.Available("motorcycle"); got != 2 {
		t.Errorf("Available(motorcycle) = %d, want 2", got)
	}
	if got := lot.Available("compact"); got != 3 {
		t.Errorf("Available(compact) = %d, want 3", got)
	}
	if got := lot.Available("large"); got != 4 {
		t.Errorf("Available(large) = %d, want 4", got)
	}
}

func TestMotorcycleTakesMotorcycleSpotFirst(t *testing.T) {
	lot := NewParkingLot(1, 1, 1)
	if lot.Park("motorcycle") <= 0 {
		t.Fatal("Park(motorcycle) failed with spots available")
	}
	if got := lot.Available("motorcycle"); got != 0 {
		t.Errorf("Available(motorcycle) = %d, want 0", got)
	}
	if got := lot.Available("compact"); got != 1 {
		t.Errorf("Available(compact) = %d, want 1 -- motorcycle must take the smallest fit", got)
	}
}

func TestMotorcycleOverflowsToCompactThenLarge(t *testing.T) {
	lot := NewParkingLot(1, 1, 1)
	lot.Park("motorcycle")
	lot.Park("motorcycle")
	if got := lot.Available("compact"); got != 0 {
		t.Errorf("Available(compact) = %d, want 0 after overflow", got)
	}
	lot.Park("motorcycle")
	if got := lot.Available("large"); got != 0 {
		t.Errorf("Available(large) = %d, want 0 after second overflow", got)
	}
}

func TestCarTakesCompactBeforeLarge(t *testing.T) {
	lot := NewParkingLot(1, 1, 1)
	if lot.Park("car") <= 0 {
		t.Fatal("Park(car) failed with spots available")
	}
	if got := lot.Available("compact"); got != 0 {
		t.Errorf("Available(compact) = %d, want 0", got)
	}
	if got := lot.Available("large"); got != 1 {
		t.Errorf("Available(large) = %d, want 1", got)
	}
}

func TestCarNeverTakesMotorcycleSpot(t *testing.T) {
	lot := NewParkingLot(5, 0, 0)
	if got := lot.Park("car"); got != -1 {
		t.Errorf("Park(car) = %d, want -1 with only motorcycle spots", got)
	}
	if got := lot.Available("motorcycle"); got != 5 {
		t.Errorf("Available(motorcycle) = %d, want 5", got)
	}
}

func TestBusNeedsALargeSpot(t *testing.T) {
	if got := NewParkingLot(5, 5, 0).Park("bus"); got != -1 {
		t.Errorf("Park(bus) = %d, want -1 with no large spots", got)
	}
	if got := NewParkingLot(0, 0, 1).Park("bus"); got <= 0 {
		t.Errorf("Park(bus) = %d, want a positive ticket", got)
	}
}

func TestParkReturnsMinusOneWhenFull(t *testing.T) {
	lot := NewParkingLot(0, 1, 0)
	if lot.Park("car") <= 0 {
		t.Fatal("first Park(car) failed")
	}
	if got := lot.Park("car"); got != -1 {
		t.Errorf("second Park(car) = %d, want -1", got)
	}
}

func TestTicketsAreUniqueAndIncreasing(t *testing.T) {
	lot := NewParkingLot(0, 3, 0)
	if t1, t2, t3 := lot.Park("car"), lot.Park("car"), lot.Park("car"); t1 != 1 || t2 != 2 || t3 != 3 {
		t.Errorf("tickets = %d,%d,%d, want 1,2,3", t1, t2, t3)
	}
}

func TestLeaveFreesTheRightSpotType(t *testing.T) {
	lot := NewParkingLot(0, 0, 1)
	ticket := lot.Park("bus")
	if got := lot.Available("large"); got != 0 {
		t.Fatalf("Available(large) = %d, want 0 while parked", got)
	}
	if !lot.Leave(ticket) {
		t.Fatal("Leave(valid ticket) = false, want true")
	}
	if got := lot.Available("large"); got != 1 {
		t.Errorf("Available(large) = %d, want 1 after Leave", got)
	}
	if lot.Park("bus") <= 0 {
		t.Error("Park(bus) failed after the spot was freed")
	}
}

func TestLeaveUnknownTicketReturnsFalse(t *testing.T) {
	if NewParkingLot(1, 1, 1).Leave(99) {
		t.Error("Leave(99) = true, want false for an unknown ticket")
	}
}

func TestLeaveTwiceReturnsFalse(t *testing.T) {
	lot := NewParkingLot(0, 1, 0)
	ticket := lot.Park("car")
	if !lot.Leave(ticket) {
		t.Fatal("first Leave = false, want true")
	}
	if lot.Leave(ticket) {
		t.Error("second Leave = true, want false")
	}
	if got := lot.Available("compact"); got != 1 {
		t.Errorf("Available(compact) = %d, want 1 -- double Leave must not double-free", got)
	}
}

func TestTicketsAreNotReusedAfterLeave(t *testing.T) {
	lot := NewParkingLot(0, 1, 0)
	t1 := lot.Park("car")
	lot.Leave(t1)
	if t2 := lot.Park("car"); t2 == t1 {
		t.Errorf("ticket %d reused, want a fresh ticket", t2)
	}
}
