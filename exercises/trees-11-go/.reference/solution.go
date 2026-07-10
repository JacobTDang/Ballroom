package main

import "math"

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// IsValidBST reports whether root is a valid binary search tree.
func IsValidBST(root *TreeNode) bool {
	var valid func(node *TreeNode, lo, hi int64) bool
	valid = func(node *TreeNode, lo, hi int64) bool {
		if node == nil {
			return true
		}
		v := int64(node.Val)
		if v <= lo || v >= hi {
			return false
		}
		return valid(node.Left, lo, v) && valid(node.Right, v, hi)
	}
	return valid(root, math.MinInt64, math.MaxInt64)
}
