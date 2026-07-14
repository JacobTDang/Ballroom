package main

import (
	"reflect"
	"testing"
)

// toLevelOrder serializes a tree to LeetCode's level-order array
// format, trimming only the trailing run of nils.
func toLevelOrder(root *TreeNode) []*int {
	if root == nil {
		return []*int{}
	}
	var out []*int
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if node == nil {
			out = append(out, nil)
			continue
		}
		v := node.Val
		out = append(out, &v)
		queue = append(queue, node.Left, node.Right)
	}
	for len(out) > 0 && out[len(out)-1] == nil {
		out = out[:len(out)-1]
	}
	return out
}

func ip(v int) *int { return &v }

func TestBuildTree(t *testing.T) {
	cases := []struct {
		preorder, inorder []int
		want              []*int
	}{
		{
			[]int{3, 9, 20, 15, 7},
			[]int{9, 3, 15, 20, 7},
			[]*int{ip(3), ip(9), ip(20), nil, nil, ip(15), ip(7)},
		},
		{[]int{-1}, []int{-1}, []*int{ip(-1)}},
		{
			[]int{1, 2, 3},
			[]int{3, 2, 1},
			[]*int{ip(1), ip(2), nil, ip(3)},
		},
		{
			[]int{1, 2, 4, 5, 3, 6, 7},
			[]int{4, 2, 5, 1, 6, 3, 7},
			[]*int{ip(1), ip(2), ip(3), ip(4), ip(5), ip(6), ip(7)},
		},
		{
			[]int{1, 2},
			[]int{1, 2},
			[]*int{ip(1), nil, ip(2)},
		},
	}

	for _, c := range cases {
		got := toLevelOrder(BuildTree(c.preorder, c.inorder))
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("BuildTree(%v, %v) level-order = %v, want %v", c.preorder, c.inorder, got, c.want)
		}
	}
}
