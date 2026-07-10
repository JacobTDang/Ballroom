# Linked List Cycle

Given `head`, the head of a linked list, determine if the linked list
has a cycle in it.

There is a cycle in a linked list if there is some node in the list
that can be reached again by continuously following the `next`
pointer.

Return `true` if there is a cycle in the linked list. Otherwise,
return `false`.

## Examples

```
Input: head = [3,2,0,-4], pos = 1
Output: true
Explanation: There is a cycle in the linked list, where the tail
connects to the 1st node (0-indexed).
```

```
Input: head = [1,2], pos = 0
Output: true
```

```
Input: head = [1], pos = -1
Output: false
Explanation: pos = -1 means there is no cycle.
```

## Constraints

- The number of nodes in the list is in the range `[0, 10^4]`.
- `-10^5 <= Node.val <= 10^5`
- `pos` is `-1` or a valid index in the linked list.
