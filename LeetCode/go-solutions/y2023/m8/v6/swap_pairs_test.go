package main

import (
	"testing"
)

func TestSwapPairs(t *testing.T) {

	t.Run("1 -> 2 -> 3 -> 4", func(t *testing.T) {
		list := &ListNode{1, &ListNode{2, &ListNode{3, &ListNode{4, nil}}}}
		newList := swapPairs(list)
		printList(newList)
	})

	t.Run("1", func(t *testing.T) {
		list := &ListNode{1, nil}
		newList := swapPairs(list)
		printList(newList)
	})

	t.Run("1 -> 2 -> 3 -> 4 -> 5", func(t *testing.T) {
		list := &ListNode{1, &ListNode{2, &ListNode{3, &ListNode{4, &ListNode{5, nil}}}}}
		newList := swapPairs(list)
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
