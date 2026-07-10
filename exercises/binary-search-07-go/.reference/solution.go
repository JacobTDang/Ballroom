package main

// FindMedianSortedArrays returns the median of the two sorted arrays
// nums1 and nums2 combined.
func FindMedianSortedArrays(nums1, nums2 []int) float64 {
	if len(nums1) > len(nums2) {
		nums1, nums2 = nums2, nums1
	}
	m, n := len(nums1), len(nums2)
	lo, hi := 0, m
	half := (m + n + 1) / 2
	const inf = 1 << 60

	for lo <= hi {
		i := lo + (hi-lo)/2
		j := half - i

		maxLeftA, minRightA := -inf, inf
		if i > 0 {
			maxLeftA = nums1[i-1]
		}
		if i < m {
			minRightA = nums1[i]
		}
		maxLeftB, minRightB := -inf, inf
		if j > 0 {
			maxLeftB = nums2[j-1]
		}
		if j < n {
			minRightB = nums2[j]
		}

		if maxLeftA <= minRightB && maxLeftB <= minRightA {
			if (m+n)%2 == 1 {
				return float64(max(maxLeftA, maxLeftB))
			}
			return float64(max(maxLeftA, maxLeftB)+min(minRightA, minRightB)) / 2.0
		} else if maxLeftA > minRightB {
			hi = i - 1
		} else {
			lo = i + 1
		}
	}
	return 0 // unreachable for valid, well-formed input
}
