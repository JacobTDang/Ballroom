package main

// ParkingLot has motorcycle, compact, and large spots. Vehicles park in
// the smallest spot type that fits and has space.
type ParkingLot struct {
	free       map[string]int
	tickets    map[int]string
	nextTicket int
}

var fits = map[string][]string{
	"motorcycle": {"motorcycle", "compact", "large"},
	"car":        {"compact", "large"},
	"bus":        {"large"},
}

func NewParkingLot(motorcycleSpots, compactSpots, largeSpots int) *ParkingLot {
	return &ParkingLot{
		free: map[string]int{
			"motorcycle": motorcycleSpots,
			"compact":    compactSpots,
			"large":      largeSpots,
		},
		tickets:    map[int]string{},
		nextTicket: 1,
	}
}

// Park parks a "motorcycle", "car", or "bus". Returns a unique positive
// ticket number, or -1 if no fitting spot is free.
func (p *ParkingLot) Park(vehicle string) int {
	for _, spot := range fits[vehicle] {
		if p.free[spot] > 0 {
			p.free[spot]--
			ticket := p.nextTicket
			p.nextTicket++
			p.tickets[ticket] = spot
			return ticket
		}
	}
	return -1
}

// Leave frees the spot held by ticket. Returns false for unknown or
// already-freed tickets.
func (p *ParkingLot) Leave(ticket int) bool {
	spot, ok := p.tickets[ticket]
	if !ok {
		return false
	}
	delete(p.tickets, ticket)
	p.free[spot]++
	return true
}

// Available returns the number of free spots of the given type.
func (p *ParkingLot) Available(spotType string) int {
	return p.free[spotType]
}
