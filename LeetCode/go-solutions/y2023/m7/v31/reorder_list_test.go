package main

import (
	"testing"
)

func TestReorderList(t *testing.T) {

	t.Run("1 -> 2 -> 3 -> 4", func(t *testing.T) {
		head := &ListNode{
			1, &ListNode{
				2, &ListNode{
					3, &ListNode{
						4, nil,
					},
				},
			},
		}

		reorderList(head)
		printList(head)

	})

	t.Run("1 -> 2 -> 3 -> 4 -> 5", func(t *testing.T) {
		head := &ListNode{
			1, &ListNode{
				2, &ListNode{
					3, &ListNode{
						4, &ListNode{5, nil},
					},
				},
			},
		}

		reorderList(head)
		printList(head)

	})

	t.Run("Empty Linked list", func(t *testing.T) {
		reorderList(nil)
	})

	t.Run("Linked list with single node", func(t *testing.T) {
		reorderList(&ListNode{0, nil})
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
