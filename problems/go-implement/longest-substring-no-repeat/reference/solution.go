package solution

// LengthOfLongestSubstring returns the length of the longest substring of s
// without repeating characters. s is ASCII; treat it as bytes.
//
// It runs in O(n) time and O(1) extra space (a fixed 256-entry table): a
// sliding window tracks the last index of each byte, and when a byte repeats
// inside the window its left edge jumps just past the previous occurrence.
func LengthOfLongestSubstring(s string) int {
	var last [256]int // last[c] = index+1 of the most recent c, 0 if unseen
	best, left := 0, 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if last[c] > left {
			left = last[c]
		}
		last[c] = i + 1
		best = max(best, i+1-left)
	}
	return best
}
