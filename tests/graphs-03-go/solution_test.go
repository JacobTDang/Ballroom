package main

import (
	"reflect"
	"sort"
	"testing"
)

// buildGraph builds nodes valued 1..len(adjList), wiring neighbors
// per adjList (1-indexed values, matching LeetCode's format), and
// returns node 1 (or nil for an empty graph).
func buildGraph(adjList [][]int) *Node {
	if len(adjList) == 0 {
		return nil
	}
	nodes := make([]*Node, len(adjList)+1) // 1-indexed
	for v := 1; v <= len(adjList); v++ {
		nodes[v] = &Node{Val: v}
	}
	for v, neighbors := range adjList {
		for _, nv := range neighbors {
			nodes[v+1].Neighbors = append(nodes[v+1].Neighbors, nodes[nv])
		}
	}
	return nodes[1]
}

// adjacency walks a graph via BFS from start and returns a map of
// val -> sorted neighbor vals, for structural comparison.
func adjacency(start *Node) map[int][]int {
	if start == nil {
		return map[int][]int{}
	}
	out := make(map[int][]int)
	visited := map[*Node]bool{start: true}
	queue := []*Node{start}
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		var nbVals []int
		for _, nb := range n.Neighbors {
			nbVals = append(nbVals, nb.Val)
			if !visited[nb] {
				visited[nb] = true
				queue = append(queue, nb)
			}
		}
		sort.Ints(nbVals)
		out[n.Val] = nbVals
	}
	return out
}

func allNodes(start *Node) []*Node {
	if start == nil {
		return nil
	}
	visited := map[*Node]bool{start: true}
	queue := []*Node{start}
	var all []*Node
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		all = append(all, n)
		for _, nb := range n.Neighbors {
			if !visited[nb] {
				visited[nb] = true
				queue = append(queue, nb)
			}
		}
	}
	return all
}

func TestCloneGraph(t *testing.T) {
	original := buildGraph([][]int{{2, 4}, {1, 3}, {2, 4}, {1, 3}})
	originalSet := make(map[*Node]bool)
	for _, n := range allNodes(original) {
		originalSet[n] = true
	}

	clone := CloneGraph(original)

	if !reflect.DeepEqual(adjacency(original), adjacency(clone)) {
		t.Errorf("clone adjacency = %v, want %v", adjacency(clone), adjacency(original))
	}

	for _, n := range allNodes(clone) {
		if originalSet[n] {
			t.Fatal("cloned node shares identity with an original node -- not a deep copy")
		}
	}
}

func TestCloneGraph_NilInput(t *testing.T) {
	if got := CloneGraph(nil); got != nil {
		t.Errorf("CloneGraph(nil) = %v, want nil", got)
	}
}

func TestCloneGraph_SingleNodeNoNeighbors(t *testing.T) {
	original := &Node{Val: 1}
	clone := CloneGraph(original)
	if clone == original {
		t.Fatal("clone shares identity with the original single node")
	}
	if clone.Val != 1 || len(clone.Neighbors) != 0 {
		t.Errorf("clone = %+v, want Val=1 with no neighbors", clone)
	}
}
