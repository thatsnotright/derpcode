# Binary Tree Level Order Traversal

Given the root of a binary tree, return its node values level by level, and again in zigzag order.

## Functions to implement

```go
type TreeNode struct {
	Val         int
	Left, Right *TreeNode
}

// LevelOrder returns the values of the tree level by level, top to bottom,
// each level left to right. A nil root yields an empty (or nil) result.
func LevelOrder(root *TreeNode) [][]int

// ZigzagLevelOrder is like LevelOrder but alternates direction: level 0 is
// left to right, level 1 is right to left, level 2 is left to right, ...
func ZigzagLevelOrder(root *TreeNode) [][]int
```

## Example

```
      3
     / \
    9  20
       / \
      15  7
```

```go
LevelOrder(root)       // [[3], [9, 20], [15, 7]]
ZigzagLevelOrder(root) // [[3], [20, 9], [15, 7]]
```

## Constraints

- The tree may be empty, a single node, or heavily skewed (depth up to 10,000, so avoid assuming shallow recursion).
- Values are arbitrary `int`s, including negatives and duplicates.
- Do not mutate the tree.

## Follow-ups

1. Return the levels bottom-up instead.
2. Return only the right-side view (the last node of each level).
3. Do it with DFS and a level parameter instead of a queue. What are the trade-offs?
