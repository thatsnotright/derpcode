package solution

import (
	"strings"
	"testing"
)

func TestLengthOfLongestSubstring(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{"example 1", "abcabcbb", 3},
		{"example 2", "bbbbb", 1},
		{"example 3", "pwwkew", 3},
		{"empty", "", 0},
		{"single char", "a", 1},
		{"all unique", "abcdef", 6},
		{"space counts", " ", 1},
		{"repeat at distance", "abba", 2},
		{"window must not move backwards", "tmmzuxt", 5},
		{"symbols and digits", "a1!a1!b2@", 6},
		{"repeat at ends", "abcda", 4},
		{"two chars alternating", "ababababab", 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LengthOfLongestSubstring(tt.s); got != tt.want {
				t.Fatalf("LengthOfLongestSubstring(%q) = %d; want %d", tt.s, got, tt.want)
			}
		})
	}
}

func TestLengthOfLongestSubstringLarge(t *testing.T) {
	// 128 distinct bytes repeated: the best window is exactly 128 wide.
	var b strings.Builder
	for b.Len() < 50000 {
		for c := 0; c < 128; c++ {
			b.WriteByte(byte(c))
		}
	}
	if got := LengthOfLongestSubstring(b.String()); got != 128 {
		t.Fatalf("repeating 128-byte alphabet: got %d; want 128", got)
	}
	if got := LengthOfLongestSubstring(strings.Repeat("a", 50000)); got != 1 {
		t.Fatalf("50000 identical bytes: got %d; want 1", got)
	}
}
