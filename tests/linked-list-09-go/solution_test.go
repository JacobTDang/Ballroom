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

func TestMergeKLists(t *testing.T) {
	cases := []struct {
		lists [][]int
		want  []int
	}{
		{[][]int{{1, 4, 5}, {1, 3, 4}, {2, 6}}, []int{1, 1, 2, 3, 4, 4, 5, 6}},
		{[][]int{}, []int{}},
		{[][]int{{}}, []int{}},
		{[][]int{{1}, {}, {2}}, []int{1, 2}},
	}

	for _, c := range cases {
		lists := make([]*ListNode, len(c.lists))
		for i, l := range c.lists {
			lists[i] = buildList(l)
		}
		got := toSlice(MergeKLists(lists))
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("MergeKLists(%v) = %v, want %v", c.lists, got, c.want)
		}
	}
}
