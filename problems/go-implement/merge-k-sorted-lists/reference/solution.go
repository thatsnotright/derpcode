package solution

import "container/heap"

// ListNode is a node in a singly linked list.
type ListNode struct {
	Val  int
	Next *ListNode
}

// nodeHeap is a min-heap of list nodes ordered by value.
type nodeHeap []*ListNode

func (h nodeHeap) Len() int           { return len(h) }
func (h nodeHeap) Less(i, j int) bool { return h[i].Val < h[j].Val }
func (h nodeHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *nodeHeap) Push(x any)        { *h = append(*h, x.(*ListNode)) }
func (h *nodeHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	old[n-1] = nil
	*h = old[:n-1]
	return x
}

// MergeKLists merges k sorted lists into one sorted list and returns its head.
// Entries in lists may be nil.
//
// It runs in O(N log k) time for N total nodes and O(k) extra space: a min-heap
// holds the current head of each non-empty list, and the smallest is repeatedly
// popped and appended to the result.
func MergeKLists(lists []*ListNode) *ListNode {
	h := make(nodeHeap, 0, len(lists))
	for _, l := range lists {
		if l != nil {
			h = append(h, l)
		}
	}
	heap.Init(&h)

	dummy := &ListNode{}
	tail := dummy
	for h.Len() > 0 {
		n := h[0]
		tail.Next = n
		tail = n
		if n.Next != nil {
			h[0] = n.Next
			heap.Fix(&h, 0)
		} else {
			heap.Pop(&h)
		}
	}
	return dummy.Next
}
