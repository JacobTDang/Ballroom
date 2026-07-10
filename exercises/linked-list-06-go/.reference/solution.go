package main

// ListNode is a singly linked list node.
type ListNode struct {
	Val  int
	Next *ListNode
}

// AddTwoNumbers adds the two numbers represented by l1 and l2 (each
// digit stored least-significant-first) and returns the sum in the
// same representation.
func AddTwoNumbers(l1, l2 *ListNode) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	carry := 0
	for l1 != nil || l2 != nil || carry != 0 {
		sum := carry
		if l1 != nil {
			sum += l1.Val
			l1 = l1.Next
		}
		if l2 != nil {
			sum += l2.Val
			l2 = l2.Next
		}
		carry = sum / 10
		cur.Next = &ListNode{Val: sum % 10}
		cur = cur.Next
	}
	return dummy.Next
}
