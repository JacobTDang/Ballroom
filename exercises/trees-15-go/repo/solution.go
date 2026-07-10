package main

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// Serialize encodes root as a string that Deserialize can turn back
// into an equivalent tree. The exact format is up to you.
func Serialize(root *TreeNode) string {
	return ""
}

// Deserialize decodes a string produced by Serialize back into the
// original tree.
func Deserialize(data string) *TreeNode {
	return nil
}
