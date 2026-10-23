package solution

import "slices"

// TreeNode is a node of a binary tree.
type TreeNode struct {
	Val         int
	Left, Right *TreeNode
}

// LevelOrder returns the values of the tree level by level, top to bottom,
// each level left to right. A nil root yields an empty (or nil) result.
func LevelOrder(root *TreeNode) [][]int {
	var levels [][]int
	if root == nil {
		return levels
	}
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		level := make([]int, 0, len(queue))
		next := make([]*TreeNode, 0, 2*len(queue))
		for _, n := range queue {
			level = append(level, n.Val)
			if n.Left != nil {
				next = append(next, n.Left)
			}
			if n.Right != nil {
				next = append(next, n.Right)
			}
		}
		levels = append(levels, level)
		queue = next
	}
	return levels
}

// ZigzagLevelOrder is like LevelOrder but alternates direction: level 0 is
// left to right, level 1 is right to left, level 2 is left to right, ...
func ZigzagLevelOrder(root *TreeNode) [][]int {
	levels := LevelOrder(root)
	for i := 1; i < len(levels); i += 2 {
		slices.Reverse(levels[i])
	}
	return levels
}
