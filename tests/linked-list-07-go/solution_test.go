package main

import "testing"

// buildCycleList builds a list from vals and, if pos >= 0, connects
// the tail's Next back to the node at pos to form a cycle.
func buildCycleList(vals []int, pos int) *ListNode {
	if len(vals) == 0 {
		return nil
	}
	nodes := make([]*ListNode, len(vals))
	for i, v := range vals {
		nodes[i] = &ListNode{Val: v}
	}
	for i := 0; i < len(nodes)-1; i++ {
		nodes[i].Next = nodes[i+1]
	}
	if pos >= 0 {
		nodes[len(nodes)-1].Next = nodes[pos]
	}
	return nodes[0]
}

func TestHasCycle(t *testing.T) {
	cases := []struct {
		vals []int
		pos  int
		want bool
	}{
		{[]int{3, 2, 0, -4}, 1, true},
		{[]int{1, 2}, 0, true},
		{[]int{1}, -1, false},
		{[]int{}, -1, false},
		{[]int{1, 2, 3}, -1, false},
		{[]int{1, 2, 3, 4, 5, 6}, 2, true},
	}

	for _, c := range cases {
		got := HasCycle(buildCycleList(c.vals, c.pos))
		if got != c.want {
			t.Errorf("HasCycle(vals=%v, pos=%d) = %v, want %v", c.vals, c.pos, got, c.want)
		}
	}
}
