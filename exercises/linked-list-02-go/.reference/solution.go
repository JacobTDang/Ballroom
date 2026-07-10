package main

// ListNode is a singly linked list node.
type ListNode struct {
	Val  int
	Next *ListNode
}

// MergeTwoLists merges two sorted linked lists into one sorted list.
func MergeTwoLists(list1, list2 *ListNode) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	for list1 != nil && list2 != nil {
		if list1.Val <= list2.Val {
			cur.Next = list1
			list1 = list1.Next
		} else {
			cur.Next = list2
			list2 = list2.Next
		}
		cur = cur.Next
	}
	if list1 != nil {
		cur.Next = list1
	} else {
		cur.Next = list2
	}
	return dummy.Next
}
