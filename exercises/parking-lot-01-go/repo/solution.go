package main

// ParkingLot has motorcycle, compact, and large spots. Vehicles park in
// the smallest spot type that fits and has space.
type ParkingLot struct {
}

func NewParkingLot(motorcycleSpots, compactSpots, largeSpots int) *ParkingLot {
	return &ParkingLot{}
}

// Park parks a "motorcycle", "car", or "bus". Returns a unique positive
// ticket number, or -1 if no fitting spot is free.
func (p *ParkingLot) Park(vehicle string) int {
	// TODO: implement
	return -1
}

// Leave frees the spot held by ticket. Returns false for unknown or
// already-freed tickets.
func (p *ParkingLot) Leave(ticket int) bool {
	// TODO: implement
	return false
}

// Available returns the number of free spots of the given type.
func (p *ParkingLot) Available(spotType string) int {
	// TODO: implement
	return 0
}
