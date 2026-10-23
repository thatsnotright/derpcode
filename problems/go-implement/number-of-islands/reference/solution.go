package solution

// NumIslands returns the number of islands in grid, where '1' is land and '0'
// is water. Cells connect horizontally and vertically.
//
// It runs in O(m*n) time and O(m*n) worst-case extra space for the explicit
// stack. Each unvisited land cell starts a new island, and an iterative flood
// fill sinks every connected land cell to '0' so it is never counted twice.
// The input grid is modified.
func NumIslands(grid [][]byte) int {
	type cell struct{ r, c int }
	dirs := [4]cell{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

	count := 0
	for r := range grid {
		for c := range grid[r] {
			if grid[r][c] != '1' {
				continue
			}
			count++
			grid[r][c] = '0'
			stack := []cell{{r, c}}
			for len(stack) > 0 {
				cur := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				for _, d := range dirs {
					nr, nc := cur.r+d.r, cur.c+d.c
					if nr < 0 || nr >= len(grid) || nc < 0 || nc >= len(grid[nr]) || grid[nr][nc] != '1' {
						continue
					}
					grid[nr][nc] = '0'
					stack = append(stack, cell{nr, nc})
				}
			}
		}
	}
	return count
}
