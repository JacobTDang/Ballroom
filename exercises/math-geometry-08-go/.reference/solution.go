package main

// DetectSquares tracks added points and counts axis-aligned squares
// formable with a query point, using a frequency map from point to the
// number of times it was added. For a query point, every previously
// added point sharing its x-coordinate forms a candidate vertical edge;
// the two horizontal partner corners at that same side length are then
// checked for existence.
type DetectSquares struct {
	points map[[2]int]int
}

func NewDetectSquares() *DetectSquares {
	return &DetectSquares{points: make(map[[2]int]int)}
}

func (d *DetectSquares) Add(point []int) {
	key := [2]int{point[0], point[1]}
	d.points[key]++
}

func (d *DetectSquares) Count(point []int) int {
	qx, qy := point[0], point[1]
	total := 0

	for p, freq := range d.points {
		if p[0] != qx || p[1] == qy {
			continue
		}
		side := p[1] - qy
		for _, cx := range [2]int{qx + side, qx - side} {
			corner1 := [2]int{cx, qy}
			corner2 := [2]int{cx, p[1]}
			f1, ok1 := d.points[corner1]
			if !ok1 {
				continue
			}
			f2, ok2 := d.points[corner2]
			if !ok2 {
				continue
			}
			total += freq * f1 * f2
		}
	}

	return total
}
