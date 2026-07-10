package main

import (
	"reflect"
	"testing"
)

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

func ip(v int) *int { return &v }

func TestRightSideView(t *testing.T) {
	cases := []struct {
		in   []*int
		want []int
	}{
		{[]*int{ip(1), ip(2), ip(3), nil, ip(5), nil, ip(4)}, []int{1, 3, 4}},
		{[]*int{ip(1), nil, ip(3)}, []int{1, 3}},
		{[]*int{}, nil},
		{[]*int{ip(1), ip(2), ip(3), ip(4)}, []int{1, 3, 4}},
	}

	for _, c := range cases {
		got := RightSideView(buildTree(c.in))
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("RightSideView(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}
