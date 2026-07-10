package main

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// MaxPathSum returns the maximum sum of any non-empty path between
// two nodes in root's tree (the path need not pass through root).
func MaxPathSum(root *TreeNode) int {
	best := root.Val
	// gain returns the best sum obtainable by extending a path
	// upward through node into exactly one of its children (or
	// neither) -- that's the only shape a parent can splice onto its
	// own path, since a path can't branch in two directions and then
	// keep going.
	var gain func(*TreeNode) int
	gain = func(node *TreeNode) int {
		if node == nil {
			return 0
		}
		leftGain := max(gain(node.Left), 0)
		rightGain := max(gain(node.Right), 0)
		if total := node.Val + leftGain + rightGain; total > best {
			best = total
		}
		return node.Val + max(leftGain, rightGain)
	}
	gain(root)
	return best
}
