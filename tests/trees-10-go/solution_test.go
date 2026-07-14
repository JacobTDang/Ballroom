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

func TestGoodNodes(t *testing.T) {
	cases := []struct {
		in   []*int
		want int
	}{
		{[]*int{ip(3), ip(1), ip(4), ip(3), nil, ip(1), ip(5)}, 4},
		{[]*int{ip(3), ip(3), nil, ip(4), ip(2)}, 3},
		{[]*int{ip(1)}, 1},
		{[]*int{ip(1), ip(2), nil, ip(3)}, 3},
		{[]*int{ip(5), ip(3), ip(3)}, 1},
	}

	for _, c := range cases {
		got := GoodNodes(buildTree(c.in))
		if got != c.want {
			t.Errorf("GoodNodes(%v) = %d, want %d", c.in, got, c.want)
		}
	}
}
