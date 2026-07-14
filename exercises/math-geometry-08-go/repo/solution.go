package main

// DetectSquares tracks added points and counts axis-aligned squares
// formable with a query point.
type DetectSquares struct {
	points map[[2]int]int
}

func NewDetectSquares() *DetectSquares {
	return &DetectSquares{points: make(map[[2]int]int)}
}

func (d *DetectSquares) Add(point []int) {
	// TODO: implement
}

func (d *DetectSquares) Count(point []int) int {
	// TODO: implement
	return 0
}
