package solution

import (
	"reflect"
	"testing"
)

// build constructs a tree from a LeetCode-style level-order slice where nil
// marks a missing child, e.g. []any{3, 9, 20, nil, nil, 15, 7}.
func build(vals []any) *TreeNode {
	if len(vals) == 0 || vals[0] == nil {
		return nil
	}
	root := &TreeNode{Val: vals[0].(int)}
	queue := []*TreeNode{root}
	i := 1
	for len(queue) > 0 && i < len(vals) {
		n := queue[0]
		queue = queue[1:]
		if i < len(vals) && vals[i] != nil {
			n.Left = &TreeNode{Val: vals[i].(int)}
			queue = append(queue, n.Left)
		}
		i++
		if i < len(vals) && vals[i] != nil {
			n.Right = &TreeNode{Val: vals[i].(int)}
			queue = append(queue, n.Right)
		}
		i++
	}
	return root
}

func normalize(got [][]int) [][]int {
	if len(got) == 0 {
		return [][]int{}
	}
	return got
}

var cases = []struct {
	name   string
	tree   []any
	level  [][]int
	zigzag [][]int
}{
	{"empty", nil, [][]int{}, [][]int{}},
	{"single", []any{1}, [][]int{{1}}, [][]int{{1}}},
	{
		"classic",
		[]any{3, 9, 20, nil, nil, 15, 7},
		[][]int{{3}, {9, 20}, {15, 7}},
		[][]int{{3}, {20, 9}, {15, 7}},
	},
	{
		"complete depth 4",
		[]any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		[][]int{{1}, {2, 3}, {4, 5, 6, 7}, {8, 9, 10, 11, 12, 13, 14, 15}},
		[][]int{{1}, {3, 2}, {4, 5, 6, 7}, {15, 14, 13, 12, 11, 10, 9, 8}},
	},
	{
		"left skewed",
		[]any{1, 2, nil, 3, nil, 4},
		[][]int{{1}, {2}, {3}, {4}},
		[][]int{{1}, {2}, {3}, {4}},
	},
	{
		"right skewed",
		[]any{1, nil, 2, nil, 3},
		[][]int{{1}, {2}, {3}},
		[][]int{{1}, {2}, {3}},
	},
	{
		"gaps and negatives",
		[]any{-1, -2, -3, nil, 4, 5, nil},
		[][]int{{-1}, {-2, -3}, {4, 5}},
		[][]int{{-1}, {-3, -2}, {4, 5}},
	},
	{
		"duplicates and zero",
		[]any{0, 0, 0, 0, nil, nil, 0},
		[][]int{{0}, {0, 0}, {0, 0}},
		[][]int{{0}, {0, 0}, {0, 0}},
	},
	{
		"one-sided children across levels",
		[]any{1, 2, 3, nil, 4, 5},
		[][]int{{1}, {2, 3}, {4, 5}},
		[][]int{{1}, {3, 2}, {4, 5}},
	},
}

func TestLevelOrder(t *testing.T) {
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := normalize(LevelOrder(build(tt.tree)))
			if !reflect.DeepEqual(got, tt.level) {
				t.Errorf("LevelOrder = %v, want %v", got, tt.level)
			}
		})
	}
}

func TestZigzagLevelOrder(t *testing.T) {
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := normalize(ZigzagLevelOrder(build(tt.tree)))
			if !reflect.DeepEqual(got, tt.zigzag) {
				t.Errorf("ZigzagLevelOrder = %v, want %v", got, tt.zigzag)
			}
		})
	}
}

func TestDeepSkewedTree(t *testing.T) {
	const depth = 10000
	root := &TreeNode{Val: 0}
	cur := root
	for i := 1; i < depth; i++ {
		cur.Left = &TreeNode{Val: i}
		cur = cur.Left
	}
	got := LevelOrder(root)
	if len(got) != depth {
		t.Fatalf("len(LevelOrder) = %d, want %d", len(got), depth)
	}
	if got[depth-1][0] != depth-1 {
		t.Errorf("last level = %v, want [%d]", got[depth-1], depth-1)
	}
}

func TestDoesNotMutateTree(t *testing.T) {
	root := build([]any{1, 2, 3, 4, 5})
	LevelOrder(root)
	ZigzagLevelOrder(root)
	want := build([]any{1, 2, 3, 4, 5})
	if !reflect.DeepEqual(root, want) {
		t.Error("tree was modified by traversal")
	}
	// Calling twice gives the same answer.
	a := ZigzagLevelOrder(root)
	b := ZigzagLevelOrder(root)
	if !reflect.DeepEqual(a, b) {
		t.Errorf("repeated calls differ: %v vs %v", a, b)
	}
}
