package main

import (
	"reflect"
	"sort"
	"testing"
)

// normalize sorts each point's own [x, y] representation-independent
// (points are fixed pairs, not reorderable) and then sorts the list
// of points lexicographically, so K Closest's any-order guarantee
// doesn't make the test order-sensitive.
func normalize(points [][]int) [][]int {
	out := make([][]int, len(points))
	copy(out, points)
	sort.Slice(out, func(i, j int) bool {
		if out[i][0] != out[j][0] {
			return out[i][0] < out[j][0]
		}
		return out[i][1] < out[j][1]
	})
	return out
}

func TestKClosest(t *testing.T) {
	cases := []struct {
		points [][]int
		k      int
		want   [][]int
	}{
		{[][]int{{1, 3}, {-2, 2}}, 1, [][]int{{-2, 2}}},
		{[][]int{{3, 3}, {5, -1}, {-2, 4}}, 2, [][]int{{3, 3}, {-2, 4}}},
		{[][]int{{0, 1}, {1, 0}}, 2, [][]int{{0, 1}, {1, 0}}},
		{[][]int{{1, 1}, {2, 2}, {3, 3}}, 3, [][]int{{1, 1}, {2, 2}, {3, 3}}},
		{[][]int{{-5, 4}, {-6, -1}, {3, 6}, {2, -2}}, 2, [][]int{{2, -2}, {-6, -1}}},
	}

	for _, c := range cases {
		got := normalize(KClosest(c.points, c.k))
		want := normalize(c.want)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("KClosest(%v, %d) = %v, want %v (order-independent)", c.points, c.k, got, want)
		}
	}
}
