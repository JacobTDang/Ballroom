package main

// FindOrder returns a valid course order satisfying every
// prerequisite pair, or an empty slice if no valid order exists.
func FindOrder(numCourses int, prerequisites [][]int) []int {
	adj := make([][]int, numCourses)
	for _, p := range prerequisites {
		course, pre := p[0], p[1]
		adj[course] = append(adj[course], pre)
	}

	const (
		unvisited = 0
		visiting  = 1
		visited   = 2
	)
	state := make([]int, numCourses)
	order := make([]int, 0, numCourses)

	var dfs func(course int) bool
	dfs = func(course int) bool {
		if state[course] == visiting {
			return false
		}
		if state[course] == visited {
			return true
		}
		state[course] = visiting
		for _, pre := range adj[course] {
			if !dfs(pre) {
				return false
			}
		}
		state[course] = visited
		order = append(order, course)
		return true
	}

	for c := 0; c < numCourses; c++ {
		if !dfs(c) {
			return []int{}
		}
	}
	return order
}
