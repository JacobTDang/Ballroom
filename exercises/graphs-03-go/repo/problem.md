# Clone Graph

Given a reference of a node in a connected undirected graph, return a
deep copy (clone) of the graph.

Each node in the graph contains a value (`int`) and a list
(`Neighbors`) of its neighbors.

## Example

```
Input: adjList = [[2,4],[1,3],[2,4],[1,3]]
Output: [[2,4],[1,3],[2,4],[1,3]]
Explanation: There are 4 nodes in the graph.
1st node (val = 1)'s neighbors are 2nd node (val = 2) and 4th node (val = 4).
2nd node (val = 2)'s neighbors are 1st node (val = 1) and 3rd node (val = 3).
3rd node (val = 3)'s neighbors are 2nd node (val = 2) and 4th node (val = 4).
4th node (val = 4)'s neighbors are 1st node (val = 1) and 3rd node (val = 3).
```

## Constraints

- The number of nodes in the graph is in the range `[0, 100]`.
- `1 <= Node.Val <= 100`
- `Node.Val` is unique for each node.
- There are no repeated edges and no self-loops in the graph.
- The graph is connected and all nodes can be visited starting from
  the given node.
