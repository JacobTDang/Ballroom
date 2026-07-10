package main

// ListNode is a singly linked list node.
type ListNode struct {
	Val  int
	Next *ListNode
}

// RemoveNthFromEnd removes the nth node from the end of head and
// returns the new head.
func RemoveNthFromEnd(head *ListNode, n int) *ListNode {
	dummy := &ListNode{Next: head}
	fast, slow := dummy, dummy
	for i := 0; i < n; i++ {
		fast = fast.Next
	}
	for fast.Next != nil {
		fast = fast.Next
		slow = slow.Next
	}
	slow.Next = slow.Next.Next
	return dummy.Next
}
