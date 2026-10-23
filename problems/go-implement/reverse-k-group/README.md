# Reverse Nodes in k-Group

Given the head of a linked list, reverse the nodes of the list `k` at a time, and return the
modified list. If the remaining nodes are fewer than `k`, leave them in their original order.

Rearrange the nodes themselves; don't just swap their values.

## Method to implement

```go
// ListNode is a node in a singly linked list.
type ListNode struct {
	Val  int
	Next *ListNode
}

// ReverseKGroup reverses the list k nodes at a time and returns the new head.
func ReverseKGroup(head *ListNode, k int) *ListNode
```

## Examples

**Example 1**

```
Input:  head = [1,2,3,4,5], k = 2
Output: [2,1,4,3,5]
```

**Example 2**

```
Input:  head = [1,2,3,4,5], k = 3
Output: [3,2,1,4,5]
```

## Constraints

- The number of nodes in the list is `n`.
- `1 <= k <= n <= 5000`
- `0 <= Node.val <= 1000`

## Follow-ups

1. Can you solve it using only `O(1)` extra memory (no recursion stack, no slice of nodes)?
2. What changes if the leftover tail (fewer than `k` nodes) should also be reversed?
3. How would you test this? Which edge cases matter most?
