# Copy List with Random Pointer

A linked list of length `n` is given such that each node contains an
additional random pointer, which could point to any node in the list,
or `null`.

Construct a deep copy of the list. The deep copy should consist of
exactly `n` brand new nodes, where each new node has its value set to
the value of its corresponding original node. Both the `next` and
`random` pointer of the new nodes should point to new nodes in the
copied list such that the pointers in the original list and copied
list represent the same list state. None of the pointers in the new
list should point to nodes in the original list.

Return the head of the copied linked list.

The linked list is represented in the input/output as a list of `n`
nodes. Each node is represented as a pair of `[val, random_index]`
where `random_index` is the index of the node (0-indexed) that the
random pointer points to, or `null` if it does not point to any node.

## Examples

```
Input: head = [[7,null],[13,0],[11,4],[10,2],[1,0]]
Output: [[7,null],[13,0],[11,4],[10,2],[1,0]]
```

```
Input: head = [[1,1],[2,1]]
Output: [[1,1],[2,1]]
```

```
Input: head = [[3,null],[3,0],[3,null]]
Output: [[3,null],[3,0],[3,null]]
```

## Constraints

- `0 <= n <= 1000`
- `-10^4 <= Node.val <= 10^4`
- `Node.random` is `null` or is pointing to some node in the linked
  list.
