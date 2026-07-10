package main

import "testing"

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

func TestIsBalanced(t *testing.T) {
	cases := []struct {
		in   []*int
		want bool
	}{
		{[]*int{ip(3), ip(9), ip(20), nil, nil, ip(15), ip(7)}, true},
		{[]*int{ip(1), ip(2), ip(2), ip(3), ip(3), nil, nil, ip(4), ip(4)}, false},
		{[]*int{}, true},
		{[]*int{ip(1)}, true},
	}

	for _, c := range cases {
		got := IsBalanced(buildTree(c.in))
		if got != c.want {
			t.Errorf("IsBalanced(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}
