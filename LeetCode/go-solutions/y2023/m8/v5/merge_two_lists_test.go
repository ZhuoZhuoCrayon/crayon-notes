package main

import (
	"testing"
)

func TestMergeTwoLists(t *testing.T) {

	t.Run("1 -> 2 -> 4, 1 -> 3 -> 4", func(t *testing.T) {
		list1 := &ListNode{1, &ListNode{2, &ListNode{4, nil}}}
		list2 := &ListNode{1, &ListNode{3, &ListNode{4, nil}}}
		newList := mergeTwoLists(list1, list2)
		printList(newList)
	})

	t.Run("nil, 0", func(t *testing.T) {
		var list1 *ListNode
		list1 = nil
		list2 := &ListNode{0, nil}
		newList := mergeTwoLists(list1, list2)
		printList(newList)
	})

	t.Run("1 -> 2 -> 4, 3 -> 5", func(t *testing.T) {
		list1 := &ListNode{1, &ListNode{2, &ListNode{4, nil}}}
		list2 := &ListNode{3, &ListNode{5, nil}}
		newList := mergeTwoLists(list1, list2)
		printList(newList)
	})
}

func assertResult(t testing.TB, got, want *ListNode) {
	t.Helper()
	if got != want {
		t.Errorf("detect error, got %v but want %v", got, want)
	}
}

func printList(head *ListNode) {
	cur := head
	for cur != nil {
		println(cur.Val)
		cur = cur.Next
	}
}
