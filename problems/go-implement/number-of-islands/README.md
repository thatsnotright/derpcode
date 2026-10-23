# Number of Islands

Given an `m x n` grid of bytes where `'1'` is land and `'0'` is water, return the number of
islands. An island is a group of land cells connected horizontally or vertically (not
diagonally). You may assume all four edges of the grid are surrounded by water.

You may modify the grid or leave it untouched; the tests accept either.

## Method to implement

```go
// NumIslands returns the number of islands in grid.
func NumIslands(grid [][]byte) int
```

## Examples

**Example 1**

```
Input:  grid = [
  "11110",
  "11010",
  "11000",
  "00000"
]
Output: 1
```

**Example 2**

```
Input:  grid = [
  "11000",
  "11000",
  "00100",
  "00011"
]
Output: 3
```

## Constraints

- `m == len(grid)`, `n == len(grid[i])`
- `0 <= m, n <= 300`
- `grid[i][j]` is `'0'` or `'1'`.

## Follow-ups

1. Can a large island overflow the stack with recursive DFS? How do you avoid that (BFS, explicit stack)?
2. Can you solve it with union-find? When is that the better choice?
3. What if cells are added one at a time and you must report the island count after each (Number of Islands II)?
