package main

import "testing"

func TestDetectCycle(t *testing.T) {

	t.Run("Cycle", func(t *testing.T) {
		tail := &ListNode{-4, nil}
		head := &ListNode{
			3, &ListNode{
				2, &ListNode{
					0, tail,
				},
			},
		}
		tail.Next = head.Next

		cycleEntryNode := detectCycle(head)
		assertResult(t, cycleEntryNode, head.Next)
	})

	t.Run("Empty Linked list", func(t *testing.T) {
		cycleEntryNode := detectCycle(nil)
		assertResult(t, cycleEntryNode, nil)
	})

	t.Run("Linked list with single node", func(t *testing.T) {
		cycleEntryNode := detectCycle(&ListNode{0, nil})
		assertResult(t, cycleEntryNode, nil)
	})
}

func assertResult(t testing.TB, got, want *ListNode) {
	t.Helper()
	if got != want {
		t.Errorf("detect error, got %v but want %v", got, want)
	}
}
