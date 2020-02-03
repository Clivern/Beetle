// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package util

import (
	"testing"
)

// TestInArray test cases
func TestInArray(t *testing.T) {
	// TestInArray
	t.Run("TestInArray", func(t *testing.T) {
		got := InArray("A", []string{"A", "B", "C", "D"})
		want := true
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	// TestInArray
	t.Run("TestInArray", func(t *testing.T) {
		got := InArray("B", []string{"A", "B", "C", "D"})
		want := true
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	// TestInArray
	t.Run("TestInArray", func(t *testing.T) {
		got := InArray("H", []string{"A", "B", "C", "D"})
		want := false
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	// TestInArray
	t.Run("TestInArray", func(t *testing.T) {
		got := InArray(1, []int{2, 3, 1})
		want := true
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	// TestInArray
	t.Run("TestInArray", func(t *testing.T) {
		got := InArray(9, []int{2, 3, 1})
		want := false
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
