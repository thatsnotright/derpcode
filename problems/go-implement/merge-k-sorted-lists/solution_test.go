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

func build(in [][]int) []*ListNode {
	lists := make([]*ListNode, len(in))
	for i, vals := range in {
		lists[i] = fromSlice(vals)
	}
	return lists
}

func TestMergeKLists(t *testing.T) {
	tests := []struct {
		name string
		in   [][]int
		want []int
	}{
		{"example 1", [][]int{{1, 4, 5}, {1, 3, 4}, {2, 6}}, []int{1, 1, 2, 3, 4, 4, 5, 6}},
		{"no lists", [][]int{}, nil},
		{"single empty list", [][]int{{}}, nil},
		{"all empty lists", [][]int{{}, {}, {}}, nil},
		{"single list", [][]int{{1, 2, 3}}, []int{1, 2, 3}},
		{"some empty", [][]int{{}, {2, 3}, {}, {1}}, []int{1, 2, 3}},
		{"negatives and duplicates", [][]int{{-3, -1, 0}, {-3, -3, 2}, {0, 0}}, []int{-3, -3, -3, -1, 0, 0, 0, 2}},
		{"disjoint ranges", [][]int{{7, 8, 9}, {1, 2, 3}, {4, 5, 6}}, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{"unequal lengths", [][]int{{1}, {2, 3, 4, 5, 6}, {0, 10}}, []int{0, 1, 2, 3, 4, 5, 6, 10}},
		{"all equal values", [][]int{{5, 5}, {5}, {5, 5, 5}}, []int{5, 5, 5, 5, 5, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			total := 0
			for _, l := range tt.in {
				total += len(l)
			}
			got := toSlice(t, MergeKLists(build(tt.in)), total)
			if !slices.Equal(got, tt.want) {
				t.Fatalf("MergeKLists(%v) = %v; want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestMergeKListsLarge(t *testing.T) {
	const k, per = 2000, 5
	in := make([][]int, k)
	var want []int
	for i := range in {
		for j := 0; j < per; j++ {
			v := (i*7919+j*104729)%20001 - 10000
			in[i] = append(in[i], v)
		}
		slices.Sort(in[i])
		want = append(want, in[i]...)
	}
	slices.Sort(want)
	got := toSlice(t, MergeKLists(build(in)), k*per)
	if !slices.Equal(got, want) {
		t.Fatalf("large merge (k=%d, %d nodes) produced wrong result", k, k*per)
	}
}
