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

// findNode locates the node with the given value, for building p/q
// arguments that are real pointers into the tree (matching the real
// LeetCode signature) rather than freestanding nodes.
func findNode(root *TreeNode, val int) *TreeNode {
	for root != nil {
		if val == root.Val {
			return root
		}
		if val < root.Val {
			root = root.Left
		} else {
			root = root.Right
		}
	}
	return nil
}

func TestLowestCommonAncestor(t *testing.T) {
	tree := []*int{ip(6), ip(2), ip(8), ip(0), ip(4), ip(7), ip(9), nil, nil, ip(3), ip(5)}
	cases := []struct {
		p, q int
		want int
	}{
		{2, 8, 6},
		{2, 4, 2},
		{0, 5, 2},
		{7, 9, 8},
	}

	for _, c := range cases {
		root := buildTree(tree)
		p := findNode(root, c.p)
		q := findNode(root, c.q)
		got := LowestCommonAncestor(root, p, q)
		if got == nil || got.Val != c.want {
			var gotVal interface{}
			if got != nil {
				gotVal = got.Val
			}
			t.Errorf("LowestCommonAncestor(p=%d, q=%d) = %v, want %d", c.p, c.q, gotVal, c.want)
		}
	}

	// small tree: root=[2,1], p=2, q=1 -> 2
	small := buildTree([]*int{ip(2), ip(1)})
	p := findNode(small, 2)
	q := findNode(small, 1)
	if got := LowestCommonAncestor(small, p, q); got == nil || got.Val != 2 {
		t.Errorf("LowestCommonAncestor(p=2, q=1) on [2,1] = %v, want 2", got)
	}
}
