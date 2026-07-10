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

func TestDiameterOfBinaryTree(t *testing.T) {
	cases := []struct {
		in   []*int
		want int
	}{
		{[]*int{ip(1), ip(2), ip(3), ip(4), ip(5)}, 3},
		{[]*int{ip(1), ip(2)}, 1},
		{[]*int{ip(1)}, 0},
		{[]*int{ip(1), ip(2), ip(3)}, 2},
	}

	for _, c := range cases {
		got := DiameterOfBinaryTree(buildTree(c.in))
		if got != c.want {
			t.Errorf("DiameterOfBinaryTree(%v) = %d, want %d", c.in, got, c.want)
		}
	}
}
