package main

import (
	"reflect"
	"testing"
)

// buildTree builds a binary tree from vals in LeetCode's level-order
// array format (nil entries are missing children).
func buildTree(vals []*int) *TreeNode {
	if len(vals) == 0 || vals[0] == nil {
		return nil
	}
	root := &TreeNode{Val: *vals[0]}
	queue := []*TreeNode{root}
	i := 1
	for len(queue) > 0 && i < len(vals) {
		node := queue[0]
		queue = queue[1:]
		if i < len(vals) {
			if vals[i] != nil {
				node.Left = &TreeNode{Val: *vals[i]}
				queue = append(queue, node.Left)
			}
			i++
		}
		if i < len(vals) {
			if vals[i] != nil {
				node.Right = &TreeNode{Val: *vals[i]}
				queue = append(queue, node.Right)
			}
			i++
		}
	}
	return root
}

// toLevelOrder serializes a tree back to the same nil-padded
// level-order format buildTree consumes, trimming only the trailing
// run of nils so results compare cleanly regardless of exact queue
// bookkeeping.
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

func TestInvertTree(t *testing.T) {
	cases := []struct {
		in   []*int
		want []*int
	}{
		{
			[]*int{ip(4), ip(2), ip(7), ip(1), ip(3), ip(6), ip(9)},
			[]*int{ip(4), ip(7), ip(2), ip(9), ip(6), ip(3), ip(1)},
		},
		{
			[]*int{ip(2), ip(1), ip(3)},
			[]*int{ip(2), ip(3), ip(1)},
		},
		{[]*int{}, []*int{}},
		{[]*int{ip(1)}, []*int{ip(1)}},
		{[]*int{ip(1), ip(2), nil}, []*int{ip(1), nil, ip(2)}},
	}

	for _, c := range cases {
		got := toLevelOrder(InvertTree(buildTree(c.in)))
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("InvertTree(%v) = %v, want %v", derefAll(c.in), derefAll(got), derefAll(c.want))
		}
	}
}

func derefAll(vals []*int) []interface{} {
	out := make([]interface{}, len(vals))
	for i, v := range vals {
		if v == nil {
			out[i] = nil
		} else {
			out[i] = *v
		}
	}
	return out
}
