package main

import "math"

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// GoodNodes counts nodes X in root's tree where no node on the path
// from root to X has a value greater than X.
func GoodNodes(root *TreeNode) int {
	var dfs func(node *TreeNode, maxSoFar int) int
	dfs = func(node *TreeNode, maxSoFar int) int {
		if node == nil {
			return 0
		}
		count := 0
		if node.Val >= maxSoFar {
			count = 1
			maxSoFar = node.Val
		}
		count += dfs(node.Left, maxSoFar)
		count += dfs(node.Right, maxSoFar)
		return count
	}
	return dfs(root, math.MinInt64)
}
