package main

// Node is an undirected graph node.
type Node struct {
	Val       int
	Neighbors []*Node
}

// CloneGraph returns a deep copy of the connected graph reachable
// from node -- every node (including neighbor references) is a
// brand new node, never shared with the input.
func CloneGraph(node *Node) *Node {
	return nil
}
