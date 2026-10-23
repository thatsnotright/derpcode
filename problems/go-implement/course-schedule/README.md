# Course Schedule II

There are `numCourses` courses labeled `0` to `numCourses - 1`. Each `prerequisites[i] = [a, b]`
means you must take course `b` before course `a`.

Return an ordering of courses that satisfies every prerequisite. If several orderings are valid,
return any of them. If it is impossible to finish all courses, return an empty slice.

## Method to implement

```go
// FindOrder returns a valid course order, or an empty slice if none exists.
func FindOrder(numCourses int, prerequisites [][]int) []int
```

## Examples

**Example 1**

```
Input:  numCourses = 2, prerequisites = [[1,0]]
Output: [0,1]
```

**Example 2**

```
Input:  numCourses = 4, prerequisites = [[1,0],[2,0],[3,1],[3,2]]
Output: [0,1,2,3]   (or [0,2,1,3])
```

**Example 3**

```
Input:  numCourses = 2, prerequisites = [[1,0],[0,1]]
Output: []
```

## Constraints

- `1 <= numCourses <= 2000`
- `0 <= len(prerequisites) <= numCourses * (numCourses - 1)`
- `prerequisites[i].length == 2`, `0 <= a, b < numCourses`, `a != b`
- All pairs `[a, b]` are distinct.

## Follow-ups

1. Can you do it with Kahn's algorithm (BFS on in-degrees) and also with DFS coloring? Compare them.
2. How do you detect and report an actual cycle, not just that one exists?
3. How would you return the order with the lexicographically smallest course numbers first?
