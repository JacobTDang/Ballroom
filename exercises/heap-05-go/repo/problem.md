# Task Scheduler

You are given an array of CPU `tasks`, each labeled with a letter
from A to Z, and a number `n`. Each CPU interval can be idle or
allow the completion of one task. Tasks can be completed in any
order, but there's a constraint: identical tasks must be separated by
at least `n` intervals due to cooldown.

Return the minimum number of CPU intervals required to complete all
tasks.

## Examples

```
Input: tasks = ["A","A","A","B","B","B"], n = 2
Output: 8
Explanation: A -> B -> idle -> A -> B -> idle -> A -> B
```

```
Input: tasks = ["A","A","A","B","B","B"], n = 0
Output: 6
Explanation: On this case any permutation of size 6 would work since
n = 0.
```

```
Input: tasks = ["A","A","A","A","A","A","B","C","D","E","F","G"], n = 2
Output: 16
Explanation: One possible solution is
A -> B -> C -> A -> D -> E -> A -> F -> G -> A -> idle -> idle -> A ->
idle -> idle -> A
```

## Constraints

- `1 <= tasks.length <= 10^4`
- `tasks[i]` is an uppercase English letter.
- `0 <= n <= 100`
