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

func TestRemoveNthFromEnd(t *testing.T) {
	cases := []struct {
		in   []int
		n    int
		want []int
	}{
		{[]int{1, 2, 3, 4, 5}, 2, []int{1, 2, 3, 5}},
		{[]int{1}, 1, []int{}},
		{[]int{1, 2}, 1, []int{1}},
		{[]int{1, 2}, 2, []int{2}},
	}

	for _, c := range cases {
		got := toSlice(RemoveNthFromEnd(buildList(c.in), c.n))
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("RemoveNthFromEnd(%v, %d) = %v, want %v", c.in, c.n, got, c.want)
		}
	}
}
