package solution

// FindOrder returns a valid course order, or an empty slice if none exists.
// prerequisites[i] = [a, b] means course b must be taken before course a.
//
// It runs in O(V+E) time and O(V+E) space using Kahn's algorithm: courses with
// no remaining prerequisites are taken first, which unlocks their dependents.
// If fewer than numCourses are taken, the rest sit on a cycle.
func FindOrder(numCourses int, prerequisites [][]int) []int {
	adj := make([][]int, numCourses)
	indeg := make([]int, numCourses)
	for _, p := range prerequisites {
		a, b := p[0], p[1]
		adj[b] = append(adj[b], a)
		indeg[a]++
	}

	order := make([]int, 0, numCourses)
	for c, d := range indeg {
		if d == 0 {
			order = append(order, c)
		}
	}
	// order doubles as the BFS queue: head advances while new courses append.
	for head := 0; head < len(order); head++ {
		for _, next := range adj[order[head]] {
			indeg[next]--
			if indeg[next] == 0 {
				order = append(order, next)
			}
		}
	}

	if len(order) != numCourses {
		return []int{}
	}
	return order
}
