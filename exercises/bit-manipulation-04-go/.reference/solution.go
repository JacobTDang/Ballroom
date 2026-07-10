package main

// ReverseBits returns n with its 32 bits in reversed order.
func ReverseBits(n uint32) uint32 {
	var result uint32
	for i := 0; i < 32; i++ {
		result = (result << 1) | (n & 1)
		n >>= 1
	}
	return result
}
