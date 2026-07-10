# Gas Station

There are `n` gas stations along a circular route, where the amount of
gas at the `i`-th station is `gas[i]`.

You have a car with an unlimited gas tank and it costs `cost[i]` of gas
to travel from the `i`-th station to its next `(i + 1)`-th station.
You begin the journey with an empty tank at one of the gas stations.

Given two integer arrays `gas` and `cost`, return the starting gas
station's index if you can travel around the circuit once in the
clockwise direction, otherwise return `-1`. If there exists a solution,
it is guaranteed to be unique.

## Example

```
Input: gas = [1,2,3,4,5], cost = [3,4,5,1,2]
Output: 3
Explanation: Starting at station 3, tank = 0 + 4 = 4.
Travel to station 4: tank = 4 - 1 + 5 = 8.
Travel to station 0: tank = 8 - 2 + 1 = 7.
Travel to station 1: tank = 7 - 3 + 2 = 6.
Travel to station 2: tank = 6 - 4 + 3 = 5.
You reach station 3 again with 5 gas remaining, having completed the
circuit.
```

## Constraints

- `n == gas.length == cost.length`
- `1 <= n <= 10^5`
- `0 <= gas[i], cost[i] <= 10^4`
