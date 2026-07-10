# Car Fleet

There are `n` cars going to the same destination along a one-lane
road. The destination is `target` miles away.

You are given two integer arrays `position` and `speed`, both of
length `n`, where `position[i]` is the position of the `i`th car and
`speed[i]` is the speed of the `i`th car (in miles per hour).

A car can never pass another car ahead of it, but it can catch up and
then travel next to it at the speed of the slower car.

A car fleet is a car or cars driving next to each other. The speed of
the fleet is the minimum speed of any car in the fleet.

If a car catches up to a car fleet right at the destination point, it
is considered to be part of the fleet.

Return the number of car fleets that will arrive at the destination.

## Examples

```
Input: target = 12, position = [10,8,0,5,3], speed = [2,4,1,1,3]
Output: 3
Explanation:
The cars starting at 10 (speed 2) and 8 (speed 4) become a fleet,
meeting each other at 12. The car starting at 0 (speed 1) does not
catch up to any other car, so it is a fleet by itself. The cars
starting at 5 (speed 1) and 3 (speed 3) become a fleet, meeting each
other at 6. There are 3 fleets in total: {10,8}, {5,3}, and {0}.
```

```
Input: target = 10, position = [3], speed = [3]
Output: 1
```

## Constraints

- `n == position.length == speed.length`
- `1 <= n <= 10^5`
- `0 < target <= 10^6`
- `0 <= position[i] < target`
- All values of `position` are unique.
- `0 < speed[i] <= 10^6`
