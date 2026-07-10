package main

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// LowestCommonAncestor returns the lowest node in the BST rooted at
// root that has both p and q as descendants (a node counts as its
// own descendant).
func LowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
	cur := root
	for cur != nil {
		switch {
		case p.Val < cur.Val && q.Val < cur.Val:
			cur = cur.Left
		case p.Val > cur.Val && q.Val > cur.Val:
			cur = cur.Right
		default:
			return cur
		}
	}
	return nil
}
