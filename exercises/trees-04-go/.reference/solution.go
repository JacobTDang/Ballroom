package main

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// IsBalanced reports whether every node's left and right subtrees
// differ in height by no more than 1.
func IsBalanced(root *TreeNode) bool {
	// height returns -1 as a sentinel meaning "already found an
	// imbalance somewhere below", short-circuiting the rest of the walk.
	var height func(*TreeNode) int
	height = func(n *TreeNode) int {
		if n == nil {
			return 0
		}
		l := height(n.Left)
		if l == -1 {
			return -1
		}
		r := height(n.Right)
		if r == -1 {
			return -1
		}
		if l-r > 1 || r-l > 1 {
			return -1
		}
		if l > r {
			return l + 1
		}
		return r + 1
	}
	return height(root) != -1
}
