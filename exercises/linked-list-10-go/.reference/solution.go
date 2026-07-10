package main

// ListNode is a singly linked list node.
type ListNode struct {
	Val  int
	Next *ListNode
}

// ReverseKGroup reverses head k nodes at a time, leaving any
// remaining group shorter than k untouched.
func ReverseKGroup(head *ListNode, k int) *ListNode {
	node := head
	count := 0
	for node != nil && count < k {
		node = node.Next
		count++
	}
	if count < k {
		return head
	}

	// node is now the head of the rest of the list, already
	// recursively reversed in groups of k.
	newHead := ReverseKGroup(node, k)

	cur := head
	prev := newHead
	for i := 0; i < k; i++ {
		next := cur.Next
		cur.Next = prev
		prev = cur
		cur = next
	}
	return prev
}
