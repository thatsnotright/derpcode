package solution

import (
	"slices"
	"testing"
)

func fromSlice(vals []int) *ListNode {
	dummy := &ListNode{}
	tail := dummy
	for _, v := range vals {
		tail.Next = &ListNode{Val: v}
		tail = tail.Next
	}
	return dummy.Next
}

// toSlice walks the list, failing the test on a cycle instead of hanging.
func toSlice(t *testing.T, head *ListNode, limit int) []int {
	t.Helper()
	var out []int
	for n := head; n != nil; n = n.Next {
		if len(out) > limit {
			t.Fatalf("list has more than %d nodes; is there a cycle?", limit)
		}
		out = append(out, n.Val)
	}
	return out
}

func seq(n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = i + 1
	}
	return out
}

func TestReverseKGroupExamples(t *testing.T) {
	tests := []struct {
		name string
		head []int
		k    int
		want []int
	}{
		{"example 1", []int{1, 2, 3, 4, 5}, 2, []int{2, 1, 4, 3, 5}},
		{"example 2", []int{1, 2, 3, 4, 5}, 3, []int{3, 2, 1, 4, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toSlice(t, ReverseKGroup(fromSlice(tt.head), tt.k), len(tt.head))
			if !slices.Equal(got, tt.want) {
				t.Fatalf("ReverseKGroup(%v, %d) = %v; want %v", tt.head, tt.k, got, tt.want)
			}
		})
	}
}

func TestReverseKGroupEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		head []int
		k    int
		want []int
	}{
		{"k equals 1 is a no-op", []int{1, 2, 3}, 1, []int{1, 2, 3}},
		{"single node", []int{7}, 1, []int{7}},
		{"k equals n reverses everything", []int{1, 2, 3, 4}, 4, []int{4, 3, 2, 1}},
		{"exact multiple of k", []int{1, 2, 3, 4, 5, 6}, 3, []int{3, 2, 1, 6, 5, 4}},
		{"tail of one left alone", []int{1, 2, 3, 4, 5, 6, 7}, 3, []int{3, 2, 1, 6, 5, 4, 7}},
		{"tail of k-1 left alone", []int{1, 2, 3, 4, 5}, 4, []int{4, 3, 2, 1, 5}},
		{"two nodes k=2", []int{1, 2}, 2, []int{2, 1}},
		{"duplicate and zero values", []int{0, 0, 1000, 5, 5}, 2, []int{0, 0, 5, 1000, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toSlice(t, ReverseKGroup(fromSlice(tt.head), tt.k), len(tt.head))
			if !slices.Equal(got, tt.want) {
				t.Fatalf("ReverseKGroup(%v, %d) = %v; want %v", tt.head, tt.k, got, tt.want)
			}
		})
	}
}

func TestReverseKGroupMovesNodesNotValues(t *testing.T) {
	head := fromSlice([]int{1, 2, 3, 4})
	var orig []*ListNode
	for n := head; n != nil; n = n.Next {
		orig = append(orig, n)
	}
	got := ReverseKGroup(head, 2)
	want := []*ListNode{orig[1], orig[0], orig[3], orig[2]}
	i := 0
	for n := got; n != nil; n = n.Next {
		if i >= len(want) || n != want[i] {
			t.Fatalf("node %d is not the expected original node; rearrange nodes instead of copying or swapping values", i)
		}
		i++
	}
	if i != len(want) {
		t.Fatalf("got %d nodes; want %d", i, len(want))
	}
}

func TestReverseKGroupLargeList(t *testing.T) {
	const n, k = 5000, 7
	got := toSlice(t, ReverseKGroup(fromSlice(seq(n)), k), n)
	want := seq(n)
	for start := 0; start+k <= n; start += k {
		slices.Reverse(want[start : start+k])
	}
	if !slices.Equal(got, want) {
		t.Fatalf("large list (n=%d, k=%d) reversed incorrectly", n, k)
	}
}

func TestReverseKGroupEveryK(t *testing.T) {
	const n = 12
	for k := 1; k <= n; k++ {
		got := toSlice(t, ReverseKGroup(fromSlice(seq(n)), k), n)
		want := seq(n)
		for start := 0; start+k <= n; start += k {
			slices.Reverse(want[start : start+k])
		}
		if !slices.Equal(got, want) {
			t.Errorf("k=%d: got %v; want %v", k, got, want)
		}
	}
}
