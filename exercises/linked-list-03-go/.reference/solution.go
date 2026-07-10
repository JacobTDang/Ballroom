package main

// ListNode is a singly linked list node.
type ListNode struct {
	Val  int
	Next *ListNode
}

// ReorderList reorders head in place from L0->L1->...->Ln into
// L0->Ln->L1->Ln-1->...
func ReorderList(head *ListNode) {
	if head == nil || head.Next == nil {
		return
	}

	// Find the middle (slow ends at the start of the second half).
	slow, fast := head, head
	for fast.Next != nil && fast.Next.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}

	// Reverse the second half.
	second := slow.Next
	slow.Next = nil
	var prev *ListNode
	for second != nil {
		next := second.Next
		second.Next = prev
		prev = second
		second = next
	}

	// Weave the two halves together.
	first := head
	for prev != nil {
		firstNext := first.Next
		secondNext := prev.Next
		first.Next = prev
		if firstNext != nil {
			prev.Next = firstNext
		}
		first = firstNext
		prev = secondNext
	}
}
