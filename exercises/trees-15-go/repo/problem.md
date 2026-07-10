# Serialize and Deserialize Binary Tree

Serialization is the process of converting a data structure into a
sequence of bits so that it can be stored or transmitted, and later
reconstructed in the same or another computer environment.

Design an algorithm to serialize and deserialize a binary tree. There
is no restriction on how your serialization/deserialization algorithm
should work — you just need to ensure that a binary tree can be
serialized to a string, and this string can be deserialized back to
the original tree structure.

There is no requirement on the format of the string, as long as
`Deserialize(Serialize(root))` reconstructs a tree with the same
structure and values as `root`.

## Examples

```
Input: root = [1,2,3,null,null,4,5]
Output: [1,2,3,null,null,4,5]
```

```
Input: root = []
Output: []
```

## Constraints

- The number of nodes in the tree is in the range `[0, 10^4]`.
- `-1000 <= Node.val <= 1000`
