# Merge k Sorted Lists

You are given an array of `k` linked lists, each sorted in ascending order. Merge all of them
into one sorted linked list and return its head.

## Method to implement

```go
// ListNode is a node in a singly linked list.
type ListNode struct {
	Val  int
	Next *ListNode
}

// MergeKLists merges k sorted lists into one sorted list and returns its head.
func MergeKLists(lists []*ListNode) *ListNode
```

## Examples

**Example 1**

```
Input:  lists = [[1,4,5],[1,3,4],[2,6]]
Output: [1,1,2,3,4,4,5,6]
```

**Example 2**

```
Input:  lists = []
Output: []
```

**Example 3**

```
Input:  lists = [[]]
Output: []
```

## Constraints

- `k == len(lists)`, `0 <= k <= 10^4`
- `0 <= len(lists[i]) <= 500`
- `-10^4 <= Node.val <= 10^4`
- Each list is sorted in ascending order.
- The sum of all list lengths is at most `10^4` in the tests here.

## Follow-ups

1. What is the time complexity with a min-heap? With divide-and-conquer pairwise merging?
2. Can you do it with `O(1)` extra space beyond the output (reusing the input nodes)?
3. How would the approach change if the lists were too large to fit in memory (streams)?
