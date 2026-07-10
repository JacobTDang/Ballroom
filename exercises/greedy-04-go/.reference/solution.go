package main

// CanCompleteCircuit returns the starting gas station index from which
// the circuit can be completed once, or -1 if no such start exists.
func CanCompleteCircuit(gas []int, cost []int) int {
	totalGas, totalCost := 0, 0
	for i := range gas {
		totalGas += gas[i]
		totalCost += cost[i]
	}
	if totalGas < totalCost {
		return -1
	}

	tank := 0
	start := 0
	for i := range gas {
		tank += gas[i] - cost[i]
		if tank < 0 {
			tank = 0
			start = i + 1
		}
	}
	return start
}
