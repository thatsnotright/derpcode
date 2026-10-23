package solution

import "testing"

// makeGrid builds a fresh mutable grid so each call may modify it freely.
func makeGrid(rows []string) [][]byte {
	grid := make([][]byte, len(rows))
	for i, r := range rows {
		grid[i] = []byte(r)
	}
	return grid
}

func TestNumIslands(t *testing.T) {
	tests := []struct {
		name string
		grid []string
		want int
	}{
		{"example 1", []string{"11110", "11010", "11000", "00000"}, 1},
		{"example 2", []string{"11000", "11000", "00100", "00011"}, 3},
		{"empty grid", []string{}, 0},
		{"empty rows", []string{""}, 0},
		{"single land", []string{"1"}, 1},
		{"single water", []string{"0"}, 0},
		{"all water", []string{"000", "000"}, 0},
		{"all land", []string{"111", "111", "111"}, 1},
		{"diagonals do not connect", []string{"101", "010", "101"}, 5},
		{"single row", []string{"1011011"}, 3},
		{"single column", []string{"1", "1", "0", "1"}, 2},
		{"snake is one island", []string{"11111", "00001", "11111", "10000", "11111"}, 1},
		{"ring around water", []string{"111", "101", "111"}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NumIslands(makeGrid(tt.grid)); got != tt.want {
				t.Fatalf("NumIslands(%v) = %d; want %d", tt.grid, got, tt.want)
			}
		})
	}
}

func TestNumIslandsLarge(t *testing.T) {
	const n = 300
	// Checkerboard: every land cell is isolated.
	grid := make([][]byte, n)
	want := 0
	for i := range grid {
		grid[i] = make([]byte, n)
		for j := range grid[i] {
			grid[i][j] = '0' + byte((i+j+1)%2)
			if grid[i][j] == '1' {
				want++
			}
		}
	}
	if got := NumIslands(grid); got != want {
		t.Fatalf("checkerboard %dx%d: got %d; want %d", n, n, got, want)
	}
}

func TestNumIslandsLargeSingleIsland(t *testing.T) {
	// A fully land grid is one big island; guards against stack blowups and
	// quadratic rescans.
	const n = 300
	grid := make([][]byte, n)
	for i := range grid {
		row := make([]byte, n)
		for j := range row {
			row[j] = '1'
		}
		grid[i] = row
	}
	if got := NumIslands(grid); got != 1 {
		t.Fatalf("all-land %dx%d: got %d; want 1", n, n, got)
	}
}
