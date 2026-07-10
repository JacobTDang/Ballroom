package main

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// BuildTree reconstructs the unique binary tree whose preorder and
// inorder traversals are preorder and inorder.
func BuildTree(preorder []int, inorder []int) *TreeNode {
	inorderIdx := make(map[int]int, len(inorder))
	for i, v := range inorder {
		inorderIdx[v] = i
	}
	pre := 0
	var build func(inLo, inHi int) *TreeNode
	build = func(inLo, inHi int) *TreeNode {
		if inLo > inHi {
			return nil
		}
		rootVal := preorder[pre]
		pre++
		root := &TreeNode{Val: rootVal}
		mid := inorderIdx[rootVal]
		root.Left = build(inLo, mid-1)
		root.Right = build(mid+1, inHi)
		return root
	}
	return build(0, len(inorder)-1)
}
