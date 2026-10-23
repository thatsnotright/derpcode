package solution

// ListNode is a node in a singly linked list.
type ListNode struct {
	Val  int
	Next *ListNode
}

// ReverseKGroup reverses the list k nodes at a time and returns the new head.
// A trailing group with fewer than k nodes keeps its original order.
//
// It runs in O(n) time and O(1) extra space: a dummy node anchors the head,
// and each full group is reversed in place and spliced back between the
// previous group's tail and the next group's head.
func ReverseKGroup(head *ListNode, k int) *ListNode {
	if k <= 1 {
		return head
	}
	dummy := &ListNode{Next: head}
	prevTail := dummy
	for {
		// Find the k-th node of this group; stop if the group is short.
		kth := prevTail
		for i := 0; i < k && kth != nil; i++ {
			kth = kth.Next
		}
		if kth == nil {
			return dummy.Next
		}

		groupHead := prevTail.Next
		nextGroup := kth.Next

		// Reverse the group, pointing its old head at the next group.
		prev, cur := nextGroup, groupHead
		for cur != nextGroup {
			cur.Next, prev, cur = prev, cur, cur.Next
		}

		prevTail.Next = kth  // kth is the new group head
		prevTail = groupHead // the old head is now the group tail
	}
}
