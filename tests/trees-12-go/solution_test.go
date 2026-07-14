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

func TestKthSmallest(t *testing.T) {
	cases := []struct {
		in   []*int
		k    int
		want int
	}{
		{[]*int{ip(3), ip(1), ip(4), nil, ip(2)}, 1, 1},
		{[]*int{ip(3), ip(1), ip(4), nil, ip(2)}, 2, 2},
		{[]*int{ip(3), ip(1), ip(4), nil, ip(2)}, 4, 4},
		{[]*int{ip(5), ip(3), ip(6), ip(2), ip(4), nil, nil, ip(1)}, 3, 3},
		{[]*int{ip(5), ip(3), ip(6), ip(2), ip(4), nil, nil, ip(1)}, 5, 5},
		{[]*int{ip(1)}, 1, 1},
	}

	for _, c := range cases {
		got := KthSmallest(buildTree(c.in), c.k)
		if got != c.want {
			t.Errorf("KthSmallest(%v, %d) = %d, want %d", c.in, c.k, got, c.want)
		}
	}
}
