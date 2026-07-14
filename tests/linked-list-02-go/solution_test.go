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

func TestMergeTwoLists(t *testing.T) {
	cases := []struct {
		l1, l2 []int
		want   []int
	}{
		{[]int{1, 2, 4}, []int{1, 3, 4}, []int{1, 1, 2, 3, 4, 4}},
		{[]int{}, []int{}, []int{}},
		{[]int{}, []int{0}, []int{0}},
		{[]int{5}, []int{1, 2, 4}, []int{1, 2, 4, 5}},
		{[]int{1, 1, 1}, []int{1, 1, 1}, []int{1, 1, 1, 1, 1, 1}},
		{[]int{-3, -1, 2}, []int{-2, 0, 5}, []int{-3, -2, -1, 0, 2, 5}},
	}

	for _, c := range cases {
		got := toSlice(MergeTwoLists(buildList(c.l1), buildList(c.l2)))
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("MergeTwoLists(%v, %v) = %v, want %v", c.l1, c.l2, got, c.want)
		}
	}
}
