package main

type ListNode struct {
	Val  int
	Next *ListNode
}

func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	newHead := &ListNode{-1, nil}

	curNodeInNewList := newHead
	curNodeInList1 := list1
	curNodeInList2 := list2

	for curNodeInList1 != nil && curNodeInList2 != nil {
		if curNodeInList1.Val < curNodeInList2.Val {
			curNodeInNewList.Next = curNodeInList1
			curNodeInNewList = curNodeInNewList.Next
			curNodeInList1 = curNodeInList1.Next
		} else {
			curNodeInNewList.Next = curNodeInList2
			curNodeInNewList = curNodeInNewList.Next
			curNodeInList2 = curNodeInList2.Next
		}
	}

	curNode := curNodeInList1
	if curNodeInList2 != nil {
		curNode = curNodeInList2
	}

	curNodeInNewList.Next = curNode

	return newHead.Next
}
