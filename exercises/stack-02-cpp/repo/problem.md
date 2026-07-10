# Min Stack

Design a stack that supports push, pop, top, and retrieving the
minimum element in constant time.

Implement the `MinStack` class:

- `MinStack()` initializes the stack object.
- `Push(val int)` pushes the element `val` onto the stack.
- `Pop()` removes the element on the top of the stack.
- `Top() int` gets the top element of the stack.
- `GetMin() int` retrieves the minimum element in the stack.

Every method must operate in `O(1)` time.

## Example

```
Input:
["MinStack","push","push","push","getMin","pop","top","getMin"]
[[],[-2],[0],[-3],[],[],[],[]]

Output:
[null,null,null,null,-3,null,0,-2]

Explanation:
MinStack minStack = new MinStack();
minStack.push(-2);
minStack.push(0);
minStack.push(-3);
minStack.getMin(); // return -3
minStack.pop();
minStack.top();    // return 0
minStack.getMin(); // return -2
```

## Constraints

- `-2^31 <= val <= 2^31 - 1`
- Methods `Pop`, `Top` and `GetMin` are called only on a non-empty
  stack.
- At most `3 * 10^4` calls will be made to `Push`, `Pop`, `Top`, and
  `GetMin`.
