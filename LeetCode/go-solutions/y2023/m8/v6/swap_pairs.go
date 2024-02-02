package main

type ListNode struct {
	Val  int
	Next *ListNode
}

func swapPairs1(head *ListNode) *ListNode {

	if head == nil || head.Next == nil {
		return head
	}

	newList := &ListNode{-1, nil}
	slow := head
	fast := head.Next
	cur := newList

	for slow != nil && slow.Next != nil {
		next := slow.Next.Next

		cur.Next = fast
		cur = cur.Next
		cur.Next = slow
		cur = cur.Next
		cur.Next = nil

		slow = next
		if slow != nil {
			fast = slow.Next
		} else {
			fast = nil
		}
	}

	if slow != nil && fast == nil {
		cur.Next = slow
	}

	return newList.Next
}

func swapPairs2(head *ListNode) *ListNode {
	dummyHead := &ListNode{-1, head}
	cur := dummyHead

	for cur.Next != nil && cur.Next.Next != nil {
		pre := cur.Next
		next := cur.Next.Next

		pre.Next = next.Next
		next.Next = pre
		cur.Next = next

		cur = pre
	}
	return dummyHead.Next
}

func swapPairs3(head *ListNode) *ListNode {
	// 空的交换结果为空
	if head == nil {
		return head
	}

	// 单节点交换结果为单节点
	if head.Next == nil {
		return head
	}

	newHead := head.Next
	head.Next = swapPairs3(newHead.Next)
	newHead.Next = head
	return newHead

	//cur := head
	//next := head.Next
	//
	//tmp := next.Next
	//next.Next = cur
	//cur.Next = swapPairs3(tmp)

	//return next
}

func swapPairs(head *ListNode) *ListNode {
	return swapPairs3(head)
}
