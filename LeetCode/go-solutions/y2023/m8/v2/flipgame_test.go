package main

import (
	"testing"
)

func TestFlipgame(t *testing.T) {

	t.Run("list with only one elements & the same value", func(t *testing.T) {
		want := 0
		got := flipgame([]int{1}, []int{1})
		if want != got {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("list with only one elements", func(t *testing.T) {
		want := 1
		got := flipgame([]int{1}, []int{2})
		if want != got {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("Case 3", func(t *testing.T) {
		want := 2
		got := flipgame([]int{1, 1}, []int{1, 2})
		if want != got {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("Case 4", func(t *testing.T) {
		want := 2
		got := flipgame([]int{1, 2, 4, 4, 7}, []int{1, 3, 4, 1, 3})
		if want != got {
			t.Errorf("got %d, want %d", got, want)
		}
	})
}
