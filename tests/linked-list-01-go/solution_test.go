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

func TestReverseList(t *testing.T) {
	cases := []struct {
		in   []int
		want []int
	}{
		{[]int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
		{[]int{1, 2}, []int{2, 1}},
		{[]int{}, []int{}},
		{[]int{7}, []int{7}},
	}

	for _, c := range cases {
		got := toSlice(ReverseList(buildList(c.in)))
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("ReverseList(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}
