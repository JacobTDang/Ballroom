package main

import "strings"

// MultiplyStrings returns the product of num1 and num2, both
// non-negative integers given as decimal strings, as a decimal string.
func MultiplyStrings(num1, num2 string) string {
	if num1 == "0" || num2 == "0" {
		return "0"
	}

	m, n := len(num1), len(num2)
	digits := make([]int, m+n)

	for i := m - 1; i >= 0; i-- {
		d1 := int(num1[i] - '0')
		for j := n - 1; j >= 0; j-- {
			d2 := int(num2[j] - '0')
			sum := d1*d2 + digits[i+j+1]
			digits[i+j+1] = sum % 10
			digits[i+j] += sum / 10
		}
	}

	start := 0
	for start < len(digits)-1 && digits[start] == 0 {
		start++
	}

	var sb strings.Builder
	for _, d := range digits[start:] {
		sb.WriteByte(byte('0' + d))
	}
	return sb.String()
}
