package main

// ListNode Definition for singly-linked list.
type ListNode struct {
	Val  int
	Next *ListNode
}

// reorderList from https://leetcode.cn/problems/reorder-list/description
func reorderList(head *ListNode) {

	if head == nil || head.Next == nil {
		return
	}

	pre, middleNode := getMiddleNode(head)
	// 断链，避免两个平行链表相交，导致 merge 失败
	pre.Next = nil

	right := reserveList(middleNode)

	head = mergeList(head, right)
}

func getMiddleNode(head *ListNode) (*ListNode, *ListNode) {
	pre := head
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		pre = slow
		slow = slow.Next
		fast = fast.Next.Next
	}

	return pre, slow
}

func reserveList(head *ListNode) *ListNode {
	newHead := head
	cur := head.Next

	for cur != nil {
		next := cur.Next
		cur.Next = newHead
		newHead = cur
		cur = next
	}
	head.Next = nil
	return newHead
}

func mergeList(left *ListNode, right *ListNode) *ListNode {
	newHead := &ListNode{-1, nil}

	cur := newHead
	leftCur := left
	rightCur := right

	for leftCur != nil && rightCur != nil {

		cur.Next = leftCur
		cur = cur.Next
		leftCur = leftCur.Next

		cur.Next = rightCur
		cur = cur.Next
		rightCur = rightCur.Next

	}

	return newHead.Next
}
