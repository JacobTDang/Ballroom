def car_fleet(target: int, position: list[int], speed: list[int]) -> int:
    """Return the number of car fleets that will arrive at target."""
    cars = sorted(zip(position, speed), reverse=True)
    stack: list[float] = []  # arrival times of fleets found so far
    for pos, spd in cars:
        t = (target - pos) / spd
        if not stack or t > stack[-1]:
            stack.append(t)
    return len(stack)
