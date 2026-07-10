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

func TestIsSameTree(t *testing.T) {
	cases := []struct {
		p, q []*int
		want bool
	}{
		{[]*int{ip(1), ip(2), ip(3)}, []*int{ip(1), ip(2), ip(3)}, true},
		{[]*int{ip(1), ip(2)}, []*int{ip(1), nil, ip(2)}, false},
		{[]*int{ip(1), ip(2), ip(1)}, []*int{ip(1), ip(1), ip(2)}, false},
		{[]*int{}, []*int{}, true},
		{[]*int{ip(1)}, []*int{}, false},
	}

	for _, c := range cases {
		got := IsSameTree(buildTree(c.p), buildTree(c.q))
		if got != c.want {
			t.Errorf("IsSameTree(%v, %v) = %v, want %v", c.p, c.q, got, c.want)
		}
	}
}
