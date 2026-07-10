package main

// ListNode is a singly linked list node.
type ListNode struct {
	Val  int
	Next *ListNode
}

// HasCycle reports whether head's list contains a cycle.
func HasCycle(head *ListNode) bool {
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
		if slow == fast {
			return true
		}
	}
	return false
}
