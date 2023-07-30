package main

// ListNode Definition for singly-linked list.
type ListNode struct {
	Val  int
	Next *ListNode
}

func detectCycle(head *ListNode) *ListNode {
	// 空链表或单节点不可能有环，直接返回
	if head == nil || head.Next == nil || head.Next.Next == nil {
		return nil
	}

	// 假设相遇点距入环点的距离为 b，环总长为 b + c，fast 经过 n 圈后追上
	// slow -> a + b
	// fast -> a + n(b + c)
	// n = 1 -> fast 和 slow 在环中最远相距 b + c - 1，假设 slow 静止，fast 速度为 1，肯定可以在 b + c - 1 秒内追上 slow
	// 此时 slow 未走完一圈
	// a + b = a + n(b + c) -> a = (n - 1)(b + c) + c -> a = c
	commonNode := getCommonNode(head)
	if commonNode == nil {
		return nil
	}

	currentNode := head
	for commonNode != currentNode {
		commonNode = commonNode.Next
		currentNode = currentNode.Next
	}

	return currentNode
}

func getCommonNode(head *ListNode) *ListNode {
	slow := head
	fast := head

	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
		if slow == fast {
			return slow
		}
	}

	return nil
}
