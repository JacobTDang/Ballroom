package main

// Node is a linked list node with an extra Random pointer that can
// point to any node in the list, or nil.
type Node struct {
	Val    int
	Next   *Node
	Random *Node
}

// CopyRandomList returns a deep copy of head — every node (including
// Random targets) is a brand new node, never shared with the input.
func CopyRandomList(head *Node) *Node {
	if head == nil {
		return nil
	}
	copies := make(map[*Node]*Node)
	for cur := head; cur != nil; cur = cur.Next {
		copies[cur] = &Node{Val: cur.Val}
	}
	for cur := head; cur != nil; cur = cur.Next {
		copies[cur].Next = copies[cur.Next]
		copies[cur].Random = copies[cur.Random]
	}
	return copies[head]
}
