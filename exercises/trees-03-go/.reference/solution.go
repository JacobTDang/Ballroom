package main

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// DiameterOfBinaryTree returns the number of edges on the longest
// path between any two nodes in root's tree.
func DiameterOfBinaryTree(root *TreeNode) int {
	best := 0
	var height func(*TreeNode) int
	height = func(n *TreeNode) int {
		if n == nil {
			return 0
		}
		l := height(n.Left)
		r := height(n.Right)
		if l+r > best {
			best = l + r
		}
		if l > r {
			return l + 1
		}
		return r + 1
	}
	height(root)
	return best
}
