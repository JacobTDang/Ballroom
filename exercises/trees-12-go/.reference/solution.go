package main

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// KthSmallest returns the kth smallest value (1-indexed) among all
// nodes in the BST rooted at root.
func KthSmallest(root *TreeNode, k int) int {
	var stack []*TreeNode
	cur := root
	for {
		for cur != nil {
			stack = append(stack, cur)
			cur = cur.Left
		}
		cur = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		k--
		if k == 0 {
			return cur.Val
		}
		cur = cur.Right
	}
}
