package solution

import (
	"slices"
	"testing"
)

func TestTrap(t *testing.T) {
	tests := []struct {
		name   string
		height []int
		want   int
	}{
		{"example 1", []int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1}, 6},
		{"example 2", []int{4, 2, 0, 3, 2, 5}, 9},
		{"nil", nil, 0},
		{"empty", []int{}, 0},
		{"single bar", []int{5}, 0},
		{"two bars", []int{3, 4}, 0},
		{"flat", []int{2, 2, 2, 2}, 0},
		{"strictly increasing", []int{1, 2, 3, 4, 5}, 0},
		{"strictly decreasing", []int{5, 4, 3, 2, 1}, 0},
		{"simple basin", []int{3, 0, 3}, 3},
		{"all zeros", []int{0, 0, 0}, 0},
		{"deep narrow well", []int{9, 0, 0, 0, 9}, 27},
		{"unequal walls", []int{5, 0, 0, 2}, 4},
		{"multiple basins", []int{3, 0, 2, 0, 4}, 7},
		{"plateau in basin", []int{2, 0, 1, 1, 0, 2}, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Trap(tt.height); got != tt.want {
				t.Fatalf("Trap(%v) = %d; want %d", tt.height, got, tt.want)
			}
		})
	}
}

func TestTrapDoesNotModifyInput(t *testing.T) {
	h := []int{4, 2, 0, 3, 2, 5}
	Trap(h)
	if want := []int{4, 2, 0, 3, 2, 5}; !slices.Equal(h, want) {
		t.Fatalf("input modified: got %v; want %v", h, want)
	}
}

func TestTrapLarge(t *testing.T) {
	const n, wall = 100000, 100000
	// Walls at both ends over a flat floor: (n-2)*wall units, beyond int32.
	h := make([]int, n)
	h[0], h[n-1] = wall, wall
	if got, want := Trap(h), (n-2)*wall; got != want {
		t.Fatalf("big basin: got %d; want %d", got, want)
	}

	// Zig-zag 0,1,0,1,...: one unit in every interior dip.
	for i := range h {
		h[i] = i % 2
	}
	if got, want := Trap(h), n/2-1; got != want {
		t.Fatalf("zig-zag: got %d; want %d", got, want)
	}
}
