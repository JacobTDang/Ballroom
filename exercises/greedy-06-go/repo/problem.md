# Merge Triplets to Form Target Triplet

A triplet is an array of three integers. You are given a 2D integer
array `triplets`, where `triplets[i] = [ai, bi, ci]` describes the
`i`-th triplet. You are also given an integer array `target =
[x, y, z]` that describes the triplet you want to obtain.

To obtain `target`, you may apply the following operation on
`triplets` any number of times (possibly zero):

- Choose two indices (0-indexed) `i` and `j` and update
  `triplets[j]` to become
  `[max(ai, aj), max(bi, bj), max(ci, cj)]`.

  For example, if `triplets[i] = [2, 5, 3]` and
  `triplets[j] = [1, 7, 5]`, `triplets[j]` becomes
  `[max(2,1), max(5,7), max(3,5)] = [2, 7, 5]`.

Return `true` if it is possible to obtain `target` as an element of
`triplets`, or `false` otherwise.

## Example

```
Input: triplets = [[2,5,3],[1,8,4],[1,7,5]], target = [2,7,5]
Output: true
Explanation: Perform the operation on triplets[0] and triplets[2]:
[2,5,3] and [1,7,5] merge into [max(2,1), max(5,7), max(3,5)] =
[2,7,5], which is target.
```

```
Input: triplets = [[3,4,5],[4,5,6]], target = [3,2,5]
Output: false
Explanation: No matter how triplets are merged, the second value of
every triplet is always >= 4, so it can never equal target[1] = 2.
```

## Constraints

- `1 <= triplets.length <= 10^5`
- `triplets[i].length == target.length == 3`
- `1 <= ai, bi, ci, x, y, z <= 1000`
