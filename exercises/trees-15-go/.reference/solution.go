package main

import (
	"strconv"
	"strings"
)

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// Serialize encodes root as a string that Deserialize can turn back
// into an equivalent tree. The exact format is up to you.
func Serialize(root *TreeNode) string {
	var sb strings.Builder
	var walk func(*TreeNode)
	walk = func(node *TreeNode) {
		if node == nil {
			sb.WriteString("#,")
			return
		}
		sb.WriteString(strconv.Itoa(node.Val))
		sb.WriteString(",")
		walk(node.Left)
		walk(node.Right)
	}
	walk(root)
	return sb.String()
}

// Deserialize decodes a string produced by Serialize back into the
// original tree.
func Deserialize(data string) *TreeNode {
	vals := strings.Split(data, ",")
	idx := 0
	var walk func() *TreeNode
	walk = func() *TreeNode {
		if idx >= len(vals) || vals[idx] == "#" {
			idx++
			return nil
		}
		v, _ := strconv.Atoi(vals[idx])
		idx++
		node := &TreeNode{Val: v}
		node.Left = walk()
		node.Right = walk()
		return node
	}
	return walk()
}
