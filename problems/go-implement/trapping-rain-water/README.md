# Trapping Rain Water

Given `n` non-negative integers representing an elevation map where each bar has width `1`,
compute how much water it can trap after raining.

## Method to implement

```go
// Trap returns the total units of water trapped by the elevation map height.
func Trap(height []int) int
```

## Examples

**Example 1**

```
Input:  height = [0,1,0,2,1,0,1,3,2,1,2,1]
Output: 6
```

**Example 2**

```
Input:  height = [4,2,0,3,2,5]
Output: 9
```

## Constraints

- `n == len(height)`
- `0 <= n <= 2 * 10^5` (the tests here go up to `10^5`)
- `0 <= height[i] <= 10^5`

## Follow-ups

1. Can you solve it in `O(n)` time and `O(1)` extra space with two pointers?
2. How does the prefix/suffix-maximum approach compare, and the monotonic-stack approach?
3. What changes in 2D (Trapping Rain Water II)?
