package solution

import "testing"

// checkOrder validates that got is a permutation of 0..n-1 in which every
// prerequisite b appears before its dependent a.
func checkOrder(t *testing.T, n int, prereqs [][]int, got []int) {
	t.Helper()
	if len(got) != n {
		t.Fatalf("order has %d courses; want %d: %v", len(got), n, got)
	}
	pos := make([]int, n)
	for i := range pos {
		pos[i] = -1
	}
	for i, c := range got {
		if c < 0 || c >= n {
			t.Fatalf("course %d out of range [0,%d)", c, n)
		}
		if pos[c] != -1 {
			t.Fatalf("course %d appears more than once", c)
		}
		pos[c] = i
	}
	for _, p := range prereqs {
		a, b := p[0], p[1]
		if pos[b] > pos[a] {
			t.Fatalf("course %d must come before %d, but order is %v", b, a, got)
		}
	}
}

func TestFindOrder(t *testing.T) {
	tests := []struct {
		name    string
		n       int
		prereqs [][]int
		ok      bool
	}{
		{"example 1", 2, [][]int{{1, 0}}, true},
		{"example 2 diamond", 4, [][]int{{1, 0}, {2, 0}, {3, 1}, {3, 2}}, true},
		{"example 3 two-cycle", 2, [][]int{{1, 0}, {0, 1}}, false},
		{"single course", 1, nil, true},
		{"no prerequisites", 5, nil, true},
		{"chain in reverse label order", 4, [][]int{{2, 3}, {1, 2}, {0, 1}}, true},
		{"cycle among other courses", 5, [][]int{{1, 0}, {2, 3}, {3, 4}, {4, 2}}, false},
		{"three-cycle", 3, [][]int{{1, 0}, {2, 1}, {0, 2}}, false},
		{"disconnected components", 6, [][]int{{1, 0}, {3, 2}, {5, 4}}, true},
		{"many dependents of one course", 5, [][]int{{1, 0}, {2, 0}, {3, 0}, {4, 0}}, true},
		{"cycle reachable from valid prefix", 4, [][]int{{1, 0}, {2, 1}, {3, 2}, {1, 3}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindOrder(tt.n, tt.prereqs)
			if !tt.ok {
				if len(got) != 0 {
					t.Fatalf("FindOrder(%d, %v) = %v; want empty (cycle)", tt.n, tt.prereqs, got)
				}
				return
			}
			checkOrder(t, tt.n, tt.prereqs, got)
		})
	}
}

func TestFindOrderLargeChain(t *testing.T) {
	const n = 2000
	var prereqs [][]int
	for i := 1; i < n; i++ {
		prereqs = append(prereqs, []int{i, i - 1})
	}
	checkOrder(t, n, prereqs, FindOrder(n, prereqs))

	// Closing the chain into a ring makes it impossible.
	prereqs = append(prereqs, []int{0, n - 1})
	if got := FindOrder(n, prereqs); len(got) != 0 {
		t.Fatalf("ring of %d courses: got order of length %d; want empty", n, len(got))
	}
}

func TestFindOrderLargeDense(t *testing.T) {
	// Each course depends on the previous few: a DAG with many edges.
	const n, band = 1000, 20
	var prereqs [][]int
	for i := 0; i < n; i++ {
		for j := max(0, i-band); j < i; j++ {
			prereqs = append(prereqs, []int{i, j})
		}
	}
	checkOrder(t, n, prereqs, FindOrder(n, prereqs))
}
