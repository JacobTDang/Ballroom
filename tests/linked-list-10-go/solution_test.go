package main

import (
	"reflect"
	"testing"
)

func buildList(vals []int) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	for _, v := range vals {
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return dummy.Next
}

func toSlice(head *ListNode) []int {
	out := []int{}
	for n := head; n != nil; n = n.Next {
		out = append(out, n.Val)
	}
	return out
}

func TestReverseKGroup(t *testing.T) {
	cases := []struct {
		in   []int
		k    int
		want []int
	}{
		{[]int{1, 2, 3, 4, 5}, 2, []int{2, 1, 4, 3, 5}},
		{[]int{1, 2, 3, 4, 5}, 3, []int{3, 2, 1, 4, 5}},
		{[]int{1, 2, 3, 4, 5}, 1, []int{1, 2, 3, 4, 5}},
		{[]int{1, 2, 3, 4, 5, 6}, 6, []int{6, 5, 4, 3, 2, 1}},
		{[]int{1}, 1, []int{1}},
		{[]int{1, 2, 3, 4, 5}, 4, []int{4, 3, 2, 1, 5}},
		{[]int{1, 2, 3, 4, 5, 6, 7, 8}, 2, []int{2, 1, 4, 3, 6, 5, 8, 7}},
	}

	for _, c := range cases {
		got := toSlice(ReverseKGroup(buildList(c.in), c.k))
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("ReverseKGroup(%v, %d) = %v, want %v", c.in, c.k, got, c.want)
		}
	}
}
