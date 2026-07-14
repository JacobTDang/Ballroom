package main

import "testing"

func isValidOrder(numCourses int, prerequisites [][]int, order []int) bool {
	if len(order) != numCourses {
		return false
	}
	pos := make(map[int]int, numCourses)
	seen := make(map[int]bool, numCourses)
	for i, c := range order {
		if seen[c] {
			return false
		}
		seen[c] = true
		pos[c] = i
	}
	for _, p := range prerequisites {
		course, pre := p[0], p[1]
		if pos[pre] >= pos[course] {
			return false
		}
	}
	return true
}

func TestFindOrder_Valid(t *testing.T) {
	prereqs := [][]int{{1, 0}, {2, 0}, {3, 1}, {3, 2}}
	order := FindOrder(4, prereqs)
	if !isValidOrder(4, prereqs, order) {
		t.Errorf("FindOrder(4, %v) = %v, not a valid topological order", prereqs, order)
	}
}

func TestFindOrder_Cycle(t *testing.T) {
	order := FindOrder(2, [][]int{{1, 0}, {0, 1}})
	if len(order) != 0 {
		t.Errorf("FindOrder(2, cyclic) = %v, want empty", order)
	}
}

func TestFindOrder_NoPrerequisites(t *testing.T) {
	order := FindOrder(3, nil)
	if !isValidOrder(3, nil, order) {
		t.Errorf("FindOrder(3, nil) = %v, not a valid order", order)
	}
}

func TestFindOrder_LinearChain(t *testing.T) {
	prereqs := [][]int{{1, 0}, {2, 1}, {3, 2}, {4, 3}}
	order := FindOrder(5, prereqs)
	if !isValidOrder(5, prereqs, order) {
		t.Errorf("FindOrder(5, %v) = %v, not a valid topological order", prereqs, order)
	}
}
