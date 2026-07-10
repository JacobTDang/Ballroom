package main

// CanFinish reports whether all numCourses courses can be completed
// given the prerequisite pairs, i.e. whether the prerequisite graph
// has no cycle.
func CanFinish(numCourses int, prerequisites [][]int) bool {
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
		return true
	}

	for c := 0; c < numCourses; c++ {
		if !dfs(c) {
			return false
		}
	}
	return true
}
