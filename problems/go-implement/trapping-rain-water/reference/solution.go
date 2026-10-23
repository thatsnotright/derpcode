package solution

// Trap returns the total units of water trapped by the elevation map height.
//
// It runs in O(n) time and O(1) extra space with two pointers. Water above a
// bar is bounded by the lower of the tallest bars on its left and right, so
// the pointer on the shorter side is always safe to settle using only its own
// side's running maximum.
func Trap(height []int) int {
	left, right := 0, len(height)-1
	leftMax, rightMax := 0, 0
	total := 0
	for left < right {
		if height[left] < height[right] {
			leftMax = max(leftMax, height[left])
			total += leftMax - height[left]
			left++
		} else {
			rightMax = max(rightMax, height[right])
			total += rightMax - height[right]
			right--
		}
	}
	return total
}
