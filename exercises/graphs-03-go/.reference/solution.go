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
	if node == nil {
		return nil
	}
	visited := make(map[*Node]*Node)
	var dfs func(*Node) *Node
	dfs = func(n *Node) *Node {
		if c, ok := visited[n]; ok {
			return c
		}
		copyNode := &Node{Val: n.Val}
		visited[n] = copyNode
		for _, nb := range n.Neighbors {
			copyNode.Neighbors = append(copyNode.Neighbors, dfs(nb))
		}
		return copyNode
	}
	return dfs(node)
}
