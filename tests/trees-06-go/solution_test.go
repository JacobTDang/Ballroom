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

func TestIsSubtree(t *testing.T) {
	cases := []struct {
		root, sub []*int
		want      bool
	}{
		{[]*int{ip(3), ip(4), ip(5), ip(1), ip(2)}, []*int{ip(4), ip(1), ip(2)}, true},
		{
			[]*int{ip(3), ip(4), ip(5), ip(1), ip(2), nil, nil, nil, nil, ip(0)},
			[]*int{ip(4), ip(1), ip(2)},
			false,
		},
		{[]*int{ip(1), ip(1)}, []*int{ip(1)}, true},
		{[]*int{ip(1)}, []*int{ip(1)}, true},
	}

	for _, c := range cases {
		got := IsSubtree(buildTree(c.root), buildTree(c.sub))
		if got != c.want {
			t.Errorf("IsSubtree(%v, %v) = %v, want %v", c.root, c.sub, got, c.want)
		}
	}
}
