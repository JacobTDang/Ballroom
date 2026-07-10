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

func TestReorderList(t *testing.T) {
	cases := []struct {
		in   []int
		want []int
	}{
		{[]int{1, 2, 3, 4}, []int{1, 4, 2, 3}},
		{[]int{1, 2, 3, 4, 5}, []int{1, 5, 2, 4, 3}},
		{[]int{1}, []int{1}},
		{[]int{1, 2}, []int{1, 2}},
	}

	for _, c := range cases {
		head := buildList(c.in)
		ReorderList(head)
		got := toSlice(head)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("ReorderList(%v) -> %v, want %v", c.in, got, c.want)
		}
	}
}
