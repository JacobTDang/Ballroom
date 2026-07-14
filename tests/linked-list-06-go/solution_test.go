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

func TestAddTwoNumbers(t *testing.T) {
	cases := []struct {
		l1, l2 []int
		want   []int
	}{
		{[]int{2, 4, 3}, []int{5, 6, 4}, []int{7, 0, 8}},
		{[]int{0}, []int{0}, []int{0}},
		{[]int{9, 9, 9, 9, 9, 9, 9}, []int{9, 9, 9, 9}, []int{8, 9, 9, 9, 0, 0, 0, 1}},
		{[]int{5}, []int{5}, []int{0, 1}},
		{[]int{1, 8}, []int{0}, []int{1, 8}},
		{[]int{2, 4, 3}, []int{5, 6, 4, 9}, []int{7, 0, 8, 9}},
	}

	for _, c := range cases {
		got := toSlice(AddTwoNumbers(buildList(c.l1), buildList(c.l2)))
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("AddTwoNumbers(%v, %v) = %v, want %v", c.l1, c.l2, got, c.want)
		}
	}
}
