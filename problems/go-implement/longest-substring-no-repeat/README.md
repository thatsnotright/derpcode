# Longest Substring Without Repeating Characters

Given a string `s`, return the length of the longest substring that contains no repeated
characters.

The input is ASCII: treat `s` as a sequence of bytes (index it with `s[i]`, not `range`).

## Method to implement

```go
// LengthOfLongestSubstring returns the length of the longest substring of s
// without repeating characters.
func LengthOfLongestSubstring(s string) int
```

## Examples

**Example 1**

```
Input:  s = "abcabcbb"
Output: 3
Explanation: "abc" is the longest.
```

**Example 2**

```
Input:  s = "bbbbb"
Output: 1
```

**Example 3**

```
Input:  s = "pwwkew"
Output: 3
Explanation: "wke"; note "pwke" is a subsequence, not a substring.
```

## Constraints

- `0 <= len(s) <= 5 * 10^4`
- `s` consists of ASCII letters, digits, symbols and spaces (bytes `0..127`).

## Follow-ups

1. Can you do it in a single pass, `O(n)` time, without re-scanning characters?
2. What changes if `s` is arbitrary Unicode? What does "character" mean then?
3. How would you return the substring itself, not just its length?
