package main

// GetSum returns a + b without using the '+' or '-' operators.
func GetSum(aIn int, bIn int) int {
	a, b := uint32(aIn), uint32(bIn)
	for b != 0 {
		carry := (a & b) << 1
		a = a ^ b
		b = carry
	}
	return int(int32(a))
}
